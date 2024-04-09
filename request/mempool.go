package request

import (
	"container/list"
	"sync"
	"time"

	// "github.com/hyperledger-labs/mirbft/crypto"
	"github.com/hyperledger-labs/mirbft/membership"
	"github.com/hyperledger-labs/mirbft/messenger"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
	logger "github.com/rs/zerolog/log"
)

type MemPool struct {
	BucketId           int
	stableMicroblocks  *list.List
	txnList            *list.List
	microblockMap      map[util.Identifier]*MicroBlock        // 所有存储下来的mb
	pendingMicroblocks map[util.Identifier]*PendingMicroblock // pending 是指还没有收集到足够的ack  处于等待不能用于共识的mb
	ackBuffer          map[util.Identifier]map[int32][]byte   // ack buffer 记录了mb收集到了那些节点的ack签名
	stableMBs          map[util.Identifier]struct{}           //stable 是指收集到了足够ack的可以用于共识了的 mb
	bsize              int                                    // number of microblocks in a proposal
	msize              int                                    // byte size of transactions in a microblock
	memsize            int                                    // number of microblocks in mempool
	currSize           int
	threshhold         int // number of acks needed for a stable microblock
	totalTx            int64
	MBmu               sync.Mutex
	Reqmu              sync.Mutex
}



// NewMemPool creates a new mempool
func NewMemPool(BucketId int) *MemPool {
	mempool := &MemPool{
		// @TODO
		// bsize:              config.GetConfig().BSize,
		// msize:              config.GetConfig().MSize,
		// memsize:            config.GetConfig().MemSize,
		// threshhold:         config.GetConfig().Q,
		BucketId:           BucketId,
		bsize:              2,
		msize:              2000,
		memsize:            3000,
		threshhold:         2,
		stableMicroblocks:  list.New(),
		microblockMap:      make(map[util.Identifier]*MicroBlock),
		pendingMicroblocks: make(map[util.Identifier]*PendingMicroblock),
		ackBuffer:          make(map[util.Identifier]map[int32][]byte),
		stableMBs:          make(map[util.Identifier]struct{}),
		currSize:           0,
		txnList:            list.New(),
	}
	return mempool
}

// AddReq adds a transaction and returns a microblock if msize is reached
// then the contained transactions should be deleted
func (pool *MemPool) AddReq(txn *Request) (bool, *MicroBlock) {
	pool.Reqmu.Lock()
	defer pool.Reqmu.Unlock()
	if pool.RemainingTx() >= int64(pool.memsize) {
		return false, nil
	}
	if pool.RemainingMB() >= int64(pool.memsize) {
		return false, nil
	}

	pool.totalTx++
	tranSize := util.SizeOf(txn)
	totalSize := tranSize + pool.currSize

	if tranSize > pool.msize {
		return false, nil
	}
	if totalSize > pool.msize {
		//do not add the curr trans, and generate a microBlock
		//set the currSize to curr trans, since it is the only one does not add to the microblock
		var id int
		pool.currSize = tranSize
		newBlock := NewMicroblock(id, pool.makeTxnSlice())
		pool.txnList.PushBack(txn)
		newBlock.Sender = membership.OwnID
		newBlock.BucketID = pool.BucketId
		newBlock.Timestamp = time.Now()
		pool.AddMicroblock(newBlock)

		if !IsBusy {
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(newBlock),
				},
			}
			for _, nodeID := range membership.AllNodeIDs() {
				messenger.EnqueueMsg(pMsg, nodeID)
			}
		} else {
			newBlock.IsForward = true
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(newBlock),
				},
			}
			pick := pickRandomNode()
			logger.Debug().Msgf("[%v] is going to forward a mb to %v ,mb hash is %x", membership.OwnID, pick, newBlock.Hash)
			messenger.EnqueueMsg(pMsg, int32(pick))
		}
		return true, newBlock

	} else if totalSize == pool.msize {
		//add the curr trans, and generate a microBlock
		var id int
		allTxn := append(pool.makeTxnSlice(), txn)

		newBlock := NewMicroblock(id, allTxn)
		pool.currSize = 0
		newBlock.Sender = membership.OwnID
		newBlock.BucketID = pool.BucketId
		newBlock.Timestamp = time.Now()

		pool.AddMicroblock(newBlock)
		// 不需要 我们可以给自己加一份
		pMsg := &pb.ProtocolMessage{
			SenderId: membership.OwnID,
			Msg: &pb.ProtocolMessage_Microblock{
				Microblock: ToProtoMicroBlock(newBlock),
			},
		}
		if !IsBusy {
			for _, nodeID := range membership.AllNodeIDs() {
				messenger.EnqueueMsg(pMsg, nodeID)
			}
		} else {
			newBlock.IsForward = true
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(newBlock),
				},
			}
			// @TODO
			pick := pickRandomNode()
			logger.Debug().Msgf("[%v] is going to forward a mb to %v , mb hash is %x", membership.OwnID, pick, newBlock.Hash)
			messenger.EnqueueMsg(pMsg, int32(pick))
		}
		return true, newBlock

	} else {
		pool.txnList.PushBack(txn)
		pool.currSize = totalSize
		return false, nil

	}
}

// AddMicroblock adds a microblock into a FIFO queue
// return an err if the queue is full (memsize)
func (pool *MemPool) AddMicroblock(mb *MicroBlock) error {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()

	_, exists := pool.microblockMap[mb.Hash]
	if exists {
		return nil
	}

	pm := &PendingMicroblock{
		Microblock: mb,
		AckMap:     make(map[int32]struct{}),
	}
	// 自己的这一份还不能添,因为在发送ACK时并不会把自己忽略
	pool.microblockMap[mb.Hash] = mb

	//check if there are some acks of this microblock arrived before
	buffer, received := pool.ackBuffer[mb.Hash]
	if received {
		// if so, add these ack to the pendingblocks
		for id, _ := range buffer {
			pm.AckMap[id] = struct{}{}
		}
		if len(pm.AckMap) >= pool.threshhold {
			if _, exists = pool.stableMBs[mb.Hash]; !exists {
				pool.stableMicroblocks.PushBack(mb)
				pool.stableMBs[mb.Hash] = struct{}{}
				delete(pool.pendingMicroblocks, mb.Hash)
			}
		} else {
			pool.pendingMicroblocks[mb.Hash] = pm
		}
	} else {
		pool.pendingMicroblocks[mb.Hash] = pm
	}
	return nil
}

// AddAck adds an ack and push a microblock into the stableMicroblocks queue if it receives enough acks
func (pool *MemPool) AddAck(ack *Ack) {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	target, received := pool.pendingMicroblocks[ack.MicroblockID]
	//check if the ack arrives before the microblock
	if received {
		target.AckMap[ack.Receiver] = struct{}{}
		if len(target.AckMap) >= pool.threshhold {
			if _, exists := pool.stableMBs[target.Microblock.Hash]; !exists {
				pool.stableMicroblocks.PushBack(target.Microblock)
				pool.stableMBs[target.Microblock.Hash] = struct{}{}
			}
		}
	}
	//ack arrives before microblock, record the number of ack received before microblock
	//let the addMicroblock do the rest.
	_, exist := pool.ackBuffer[ack.MicroblockID]
	if exist {
		pool.ackBuffer[ack.MicroblockID][ack.Receiver] = ack.Signature
	} else {
		temp := make(map[int32][]byte, 0)
		temp[ack.Receiver] = ack.Signature
		pool.ackBuffer[ack.MicroblockID] = temp
	}

}

// GeneratePayload generates a list of microblocks according to bsize
// if the remaining microblocks is less than bsize then return all
func (pool *MemPool) GeneratePayloadWithSize(batchSize int) *Payload {
	sigMap := make(map[util.Identifier]map[int32][]byte, 0)
	microblockList := make([]*MicroBlock, 0)
	for i := 0; i < batchSize; i++ {
		mb := pool.front()
		if mb == nil {
			break
		}
		microblockList = append(microblockList, mb)

		sigs := make(map[int32][]byte, 0)
		count := 0
		for id, sig := range pool.ackBuffer[mb.Hash] {
			count++
			sigs[id] = sig
			delete(pool.ackBuffer, mb.Hash)
		}
		sigMap[mb.Hash] = sigs
	}
	return NewPayload(microblockList, sigMap) // payload 带有mb以及他们相关的ack收集信息以便向其他人证明信息质量
}

// CheckExistence checks if the referred microblocks in the proposal exists
// in the mempool and return missing ones if there's any
// return true if there's no missing transactions

// RemoveMicroblock removes reffered microblocks from the mempool
func (pool *MemPool) RemoveMicroblock(id util.Identifier) error {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	_, exists := pool.microblockMap[id]
	if exists {
		delete(pool.microblockMap, id)
	}
	_, exists = pool.stableMBs[id]
	if exists {
		delete(pool.stableMBs, id)
	}
	return nil
}

// FindMicroblock finds a reffered microblock
func (pool *MemPool) FindMicroblock(id util.Identifier) (bool, *MicroBlock) {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	mb, found := pool.microblockMap[id]
	return found, mb
}

func (pool *MemPool) CheckExistence(MBHashList [][]byte) (bool, []util.Identifier) {
	id := make([]util.Identifier, 0)
	return false, id
}

// FillProposal pulls microblocks from the mempool and build a pending block,
// a pending block should include the proposal, micorblocks that already exist,
// and a missing list if there's any
func (pool *MemPool) FillProposal(MBHashList [][]byte) *PendingBlock {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	existingBlocks := make([]*MicroBlock, 0)
	missingBlocks := make(map[util.Identifier]struct{}, 0)
	for _, id := range MBHashList {
		found := false
		_, exists := pool.pendingMicroblocks[util.BytesToIdentifier(id)]
		if exists {
			found = true
			existingBlocks = append(existingBlocks, pool.pendingMicroblocks[util.BytesToIdentifier(id)].Microblock)
			delete(pool.pendingMicroblocks, util.BytesToIdentifier(id))
		}
		for e := pool.stableMicroblocks.Front(); e != nil; e = e.Next() {
			mb := e.Value.(*MicroBlock)
			if mb.Hash == util.BytesToIdentifier(id) {
				found = true
				existingBlocks = append(existingBlocks, mb)
				pool.stableMicroblocks.Remove(e)
				break
			}
		}
		if !found {
			missingBlocks[util.BytesToIdentifier(id)] = struct{}{}
		}
	}
	return NewPendingBlock(MBHashList, missingBlocks, existingBlocks)
}


// 我们需要尽可能的移除这个线程，这是一个严重影响性能的地方
// 本质是为了防止内存池中残留Req， 也就是在CutBatchs
// @DONE
func (pool *MemPool) ForceClear() {
	pool.Reqmu.Lock()
	defer pool.Reqmu.Unlock()
	if pool.currSize > 0 {
		logger.Info().Msg("Force clear triggered.")
		var id int
		allTxn := pool.makeTxnSlice()
		newBlock := NewMicroblock(id, allTxn)
		pool.currSize = 0
		newBlock.Sender = membership.OwnID
		newBlock.BucketID = pool.BucketId
		newBlock.Timestamp = time.Now()
		// 将新的微块添加到内存池中
		pool.AddMicroblock(newBlock)
		// 将新的微块发送给其他节点
		pMsg := &pb.ProtocolMessage{
			SenderId: membership.OwnID,
			Msg: &pb.ProtocolMessage_Microblock{
				Microblock: ToProtoMicroBlock(newBlock),
			},
		}
		for _, nodeID := range membership.AllNodeIDs() {
			messenger.EnqueueMsg(pMsg, nodeID)
		}
	}
}

func (pool *MemPool) IsStable(id util.Identifier) bool {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	_, exists := pool.stableMBs[id]
	if exists {
		return true
	}
	return false
}

func (pool *MemPool) TotalTx() int64 {
	return pool.totalTx
}

func (pool *MemPool) RemainingTx() int64 {
	return int64(pool.txnList.Len())
}

func (pool *MemPool) TotalMB() int64 {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	return int64(len(pool.microblockMap))
}

func (pool *MemPool) RemainingMB() int64 {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	return int64(len(pool.pendingMicroblocks) + pool.stableMicroblocks.Len())
}

func (pool *MemPool) AckList(id util.Identifier) []int32 {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	pmb, exists := pool.pendingMicroblocks[id]
	if exists {
		nodes := make([]int32, 0, len(pmb.AckMap))
		for k, _ := range pmb.AckMap {
			nodes = append(nodes, k)
		}
		return nodes
	}
	return nil
}

func (pool *MemPool) front() *MicroBlock {

	if pool.stableMicroblocks.Len() == 0 {
		return nil
	}
	ele := pool.stableMicroblocks.Front()
	val, ok := ele.Value.(*MicroBlock)
	if !ok {
		return nil
	}
	pool.stableMicroblocks.Remove(ele)
	return val
}

func (pool *MemPool) makeTxnSlice() []*Request {
	allTxn := make([]*Request, 0)
	for pool.txnList.Len() > 0 {
		e := pool.txnList.Front()
		allTxn = append(allTxn, e.Value.(*Request))
		pool.txnList.Remove(e)
	}
	return allTxn
}
