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
	MBmu                 sync.Mutex
	Reqmu                 sync.Mutex
	lastReqTime        time.Time
	autoClearThreshold time.Duration
}

func (pool *MemPool) AutoClear() {
	// 计算当前时间与上次添加请求时间的间隔
	for {			
		pool.Reqmu.Lock()
		duration := time.Since(pool.lastReqTime)
		// 如果间隔超过自动清理阈值，则执行清理操作
		if duration > pool.autoClearThreshold && pool.currSize > 0 {
			logger.Info().Msgf("Duration: %d",duration.Nanoseconds())
			logger.Info().Msg("Auto clear triggered.")

			// 创建一个新的微块
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
			pool.lastReqTime = time.Now()
			
		}
		pool.Reqmu.Unlock()
	}
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
		lastReqTime:        time.Now(),
		autoClearThreshold: 10000000000,
	}
	// go mempool.AutoClear()
	return mempool
}

// AddReq adds a transaction and returns a microblock if msize is reached
// then the contained transactions should be deleted
func (pool *MemPool) AddReq(txn *Request) (bool, *MicroBlock) {
	// if txn.Msg.RequestId.ClientSn >= 3800 {
	// 	logger.Info().Msgf("Txn in mempool %d at bucket %d", txn.Msg.RequestId.ClientSn, pool.BucketId)
	// 	logger.Info().Msgf("This pool size is %d,req:%d", pool.currSize, pool.txnList.Len())
	// }
	pool.Reqmu.Lock()
	defer pool.Reqmu.Unlock()
	pool.lastReqTime = time.Now()
	if pool.RemainingTx() >= int64(pool.memsize) {
		//log.Warningf("mempool's tx list is full")
		return false, nil
	}
	if pool.RemainingMB() >= int64(pool.memsize) {
		//log.Warningf("mempool's mb is full")
		return false, nil
	}

	pool.totalTx++

	// get the size of the structure. txn is the pointer.
	tranSize := util.SizeOf(txn)
	totalSize := tranSize + pool.currSize
	// logger.Printf("Txn size is %d ,cur mempool size is %d\n, cutting size is %d", tranSize, totalSize, pool.msize)

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
		pMsg := &pb.ProtocolMessage{
			SenderId: membership.OwnID,
			Msg: &pb.ProtocolMessage_Microblock{
				Microblock: ToProtoMicroBlock(newBlock),
			},
		}
		for _, nodeID := range membership.AllNodeIDs() {
			messenger.EnqueueMsg(pMsg, nodeID)
		}
		// logger.Info().Msgf("Generating a block %x by %d",newBlock.Hash,newBlock.Sender)
		return true, newBlock

	} else if (totalSize == pool.msize || txn.Msg.RequestId.ClientSn >= 3700){
		//add the curr trans, and generate a microBlock
		var id int
		allTxn := append(pool.makeTxnSlice(), txn)

		newBlock := NewMicroblock(id, allTxn)
		pool.currSize = 0
		newBlock.Sender = membership.OwnID
		newBlock.BucketID = pool.BucketId
		newBlock.Timestamp = time.Now()

		pool.AddMicroblock(newBlock)

		// ack := &Ack{
		// 	Receiver: membership.OwnID,
		// 	BucketID: pool.BucketId,
		// 	MicroblockID: newBlock.Hash,
		// 	Signature: nil,
		// }
		// hash := crypto.Hash(ack.AckDigest())
		// signature, err := crypto.Sign(hash, sk)

		// pool.AddAck()
		pMsg := &pb.ProtocolMessage{
			SenderId: membership.OwnID,
			Msg: &pb.ProtocolMessage_Microblock{
				Microblock: ToProtoMicroBlock(newBlock),
			},
		}
		for _, nodeID := range membership.AllNodeIDs() {
			messenger.EnqueueMsg(pMsg, nodeID)
		}
		// logger.Info().Msgf("Generating a block %x by %d",newBlock.Hash,newBlock.Sender)
		return true, newBlock

	} else {
		pool.txnList.PushBack(txn)
		pool.currSize = totalSize
		// logger.Info().Msg("Still not return yet")
		return false, nil

	}
}

// AddMicroblock adds a microblock into a FIFO queue
// return an err if the queue is full (memsize)
func (pool *MemPool) AddMicroblock(mb *MicroBlock) error {
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	// logger.Info().Any("MS:%d",len(ReceivedMBs)).Msg("About to add mb")
	//if pool.microblocks.Len() >= pool.memsize {
	//	return errors.New("the memory queue is full")
	//}
	_, exists := pool.microblockMap[mb.Hash]
	if exists {
		return nil
	}

	pm := &PendingMicroblock{
		Microblock: mb,
		AckMap:     make(map[int32]struct{}),
	}
	// 自己的这一份还不能添
	// pm.AckMap[mb.Sender] = struct{}{}
	pool.microblockMap[mb.Hash] = mb

	//check if there are some acks of this microblock arrived before
	buffer, received := pool.ackBuffer[mb.Hash]
	if received {
		// if so, add these ack to the pendingblocks
		for id, _ := range buffer {
			//pool.pendingMicroblocks[mb.Hash].ackMap[ack] = struct{}{}
			pm.AckMap[id] = struct{}{}
		}
		if len(pm.AckMap) >= pool.threshhold {
			// logger.Info().Msgf("Stable %x at pool:%d",mb.Hash,pool.BucketId)
			if _, exists = pool.stableMBs[mb.Hash]; !exists {
				pool.stableMicroblocks.PushBack(mb)
				pool.stableMBs[mb.Hash] = struct{}{}
				delete(pool.pendingMicroblocks, mb.Hash)
				//log.Debugf("microblock id: %x becomes stable from buffer", mb.Hash)
			}
		} else {
			// logger.Info().Msgf("Pending %x at pool:%d",mb.Hash,pool.BucketId)
			pool.pendingMicroblocks[mb.Hash] = pm
		}
	} else {
		// logger.Info().Msgf("Pending %x at pool:%d",mb.Hash,pool.BucketId)
		pool.pendingMicroblocks[mb.Hash] = pm
	}
	return nil
}

// AddAck adds an ack and push a microblock into the stableMicroblocks queue if it receives enough acks
func (pool *MemPool) AddAck(ack *Ack) {
	// logger.Info().Any("Confirm MB:",ack.MicroblockID).Any("Confirm Peer:",ack.Receiver).Any("Confirm Bucket:",ack.BucketID).Msg("Adding ack now!")
	pool.MBmu.Lock()
	defer pool.MBmu.Unlock()
	target, received := pool.pendingMicroblocks[ack.MicroblockID]
	//check if the ack arrives before the microblock
	if received {
		target.AckMap[ack.Receiver] = struct{}{}
		if len(target.AckMap) >= pool.threshhold {
			if _, exists := pool.stableMBs[target.Microblock.Hash]; !exists {
				// logger.Info().Msgf("Moving pending to stable: %x at pool %d",target.Microblock.Hash,pool.BucketId)
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
		//log.Debugf("microblock id: %x is deleted from mempool when proposing", mb.Hash)
		microblockList = append(microblockList, mb)

		sigs := make(map[int32][]byte, 0)
		count := 0
		for id, sig := range pool.ackBuffer[mb.Hash] {
			count++
			sigs[id] = sig
			// @TODO
			// if count == config.Configuration.Q {
			// 	break
			// }
		}
		sigMap[mb.Hash] = sigs
	}
	// logger.Info().Msgf("Generating a mb with sig:",len(sigMap))
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
			//log.Debugf("microblock id: %x is deleted from pending when filling", id)
		}
		// logger.Info().Msgf("Searching stable ... len: %d",pool.stableMicroblocks.Len())
		for e := pool.stableMicroblocks.Front(); e != nil; e = e.Next() {
			// do something with e.Value
			mb := e.Value.(*MicroBlock)
			// logger.Info().Msgf("%x vs %x",mb.Hash,util.BytesToIdentifier(id))
			if mb.Hash == util.BytesToIdentifier(id) {
				// logger.Info().Msgf("Found mb in stable ... len: %x",mb.Hash)
				existingBlocks = append(existingBlocks, mb)
				found = true
				pool.stableMicroblocks.Remove(e)
				//log.Debugf("microblock id: %x is deleted from stable when filling", mb.Hash)
				break
			}
		}
		if !found {
			missingBlocks[util.BytesToIdentifier(id)] = struct{}{}
		}
	}
	return NewPendingBlock(MBHashList, missingBlocks, existingBlocks)
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
	// pool.RemoveMicroblock(ele.(*util.Identifier))
	pool.stableMicroblocks.Remove(ele)
	// logger.Info().Msgf("Popping out a mb %x with sig: %d",val.Hash, len(pool.ackBuffer[val.Hash]))
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
