package request

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger-labs/mirbft/config"
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
	msize              int                                    // byte size of transactions in a microblock
	currSize           int
	threshhold         int // number of acks needed for a stable microblock
	MBmu               sync.Mutex
	ReqMutex           sync.Mutex
	// 统计量
}

// NewMemPool creates a new mempool
func NewMemPool(BucketId int) *MemPool {
	mempool := &MemPool{
		msize:              config.Config.PoolMSize,
		threshhold:         config.Config.PoolStableThreshhold,
		BucketId:           BucketId,
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

	pool.ReqMutex.Lock()
	defer pool.ReqMutex.Unlock()

	pool.currSize = pool.currSize + 1

	if pool.currSize >= pool.msize {
		//do not add the curr trans, and generate a microBlock
		//set the currSize to curr trans, since it is the only one does not add to the microblock
		var id int
		pool.currSize = 0
		pool.txnList.PushBack(txn)
		newMicroBlock := NewMicroblock(id, pool.makeTxnSlice())
		newMicroBlock.Sender = membership.OwnID
		newMicroBlock.BucketID = pool.BucketId
		HandleMicroblock(newMicroBlock)
		txnsInfo := ""
		for _, txn := range newMicroBlock.Txns {
			txnsInfo += fmt.Sprintf("(%v:%v) ", txn.Msg.RequestId.ClientId, txn.Msg.RequestId.ClientSn)
		}
		logger.Info().Msgf("Generate at bucket  %d and peer %d microblock %x containing Req:%s", pool.BucketId, membership.OwnID, newMicroBlock.Hash, txnsInfo)

		if !CheckLoadStatus() {
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(newMicroBlock),
				},
			}
			for _, nodeID := range membership.AllNodeIDs() {
				sampleTime := time.Now().Unix()
				UpdateLoadStatus(sampleTime)
				if nodeID == membership.OwnID {
					continue
				}
				messenger.EnqueueMsg(pMsg, nodeID)
			}
		} else {
			newMicroBlock.IsForward = true
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(newMicroBlock),
				},
			}
			pick := pickRandomNode()
			logger.Info().Msgf("Peer %v is going to forward a mb to %v ,mb hash is %x", membership.OwnID, pick, newMicroBlock.Hash)
			messenger.EnqueueMsg(pMsg, int32(pick))
		}
		return true, newMicroBlock

	} else {
		pool.txnList.PushBack(txn)
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

// @TODO ForceClear 会导致严重的带宽消耗，每次ForceClear都会导致为一个Req的长度的MB进行广播和共识..
// 我们必须谨慎的进行ForceClear，并尽可能的积攒足够的Req在打包和成MB

func (pool *MemPool) ForceClear() bool {
	pool.ReqMutex.Lock()
	defer pool.ReqMutex.Unlock()
	if pool.currSize > 0 {
		var id int
		allTxn := pool.makeTxnSlice()
		newBlock := NewMicroblock(id, allTxn)
		pool.currSize = 0
		newBlock.Sender = membership.OwnID
		newBlock.BucketID = pool.BucketId
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
			if nodeID == membership.OwnID {
				HandleMicroblock(newBlock)
				continue
			}
			messenger.EnqueueMsg(pMsg, nodeID)
		}
		return true
	}
	return false

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
