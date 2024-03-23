package mempool

import (
	"container/list"
	"github.com/hyperledger-labs/mirbft/mypart/microblock"
	"github.com/hyperledger-labs/mirbft/request"
	"github.com/hyperledger-labs/mirbft/util"
	logger "github.com/rs/zerolog/log"
	"sync"
)

type MemPool struct {
	stableMicroblocks  *list.List
	txnList            *list.List
	microblockMap      map[util.Identifier]*microblock.MicroBlock         // 所有存储下来的mb
	pendingMicroblocks map[util.Identifier]*microblock.PendingMicroblock  // pending 是指还没有收集到足够的ack  处于等待不能用于共识的mb
	ackBuffer          map[util.Identifier]map[util.NodeID]util.Signature // ack buffer 记录了mb收集到了那些节点的ack签名
	stableMBs          map[util.Identifier]struct{}                       //stable 是指收集到了足够ack的可以用于共识了的 mb
	bsize              int                                                // number of microblocks in a proposal
	msize              int                                                // byte size of transactions in a microblock
	memsize            int                                                // number of microblocks in mempool
	currSize           int
	threshhold         int // number of acks needed for a stable microblock
	totalTx            int64
	mu                 sync.Mutex
}

// NewMemPool creates a new naive mempool
func NewMemPool() *MemPool {
	return &MemPool{
		// @TODO
		// bsize:              config.GetConfig().BSize,
		// msize:              config.GetConfig().MSize,
		// memsize:            config.GetConfig().MemSize,
		// threshhold:         config.GetConfig().Q,
		bsize:              2,
		msize:              1000,
		memsize:            5,
		threshhold:         2,
		stableMicroblocks:  list.New(),
		microblockMap:      make(map[util.Identifier]*microblock.MicroBlock),
		pendingMicroblocks: make(map[util.Identifier]*microblock.PendingMicroblock),
		ackBuffer:          make(map[util.Identifier]map[util.NodeID]util.Signature),
		stableMBs:          make(map[util.Identifier]struct{}),
		currSize:           0,
		txnList:            list.New(),
	}
}

// AddTxn adds a transaction and returns a microblock if msize is reached
// then the contained transactions should be deleted
func (am *MemPool) AddTxn(txn *request.Request) (bool, *microblock.MicroBlock) {
	// mempool is full
	// log.Printf("Txn in mempool %x\n",txn.Digest)
	if am.RemainingTx() >= int64(am.memsize) {
		//log.Warningf("mempool's tx list is full")
		return false, nil
	}
	if am.RemainingMB() >= int64(am.memsize) {
		//log.Warningf("mempool's mb is full")
		return false, nil
	}
	am.totalTx++

	// get the size of the structure. txn is the pointer.
	tranSize := util.SizeOf(txn)
	totalSize := tranSize + am.currSize
	// logger.Printf("Txn size is %d ,cur mempool size is %d\n, cutting size is %d", tranSize, totalSize, am.msize)

	if tranSize > am.msize {
		return false, nil
	}
	if totalSize > am.msize {
		//do not add the curr trans, and generate a microBlock
		//set the currSize to curr trans, since it is the only one does not add to the microblock
		var id util.Identifier
		am.currSize = tranSize
		newBlock := microblock.NewMicroblock(id, am.makeTxnSlice())
		am.txnList.PushBack(txn)
		return true, newBlock

	} else if totalSize == am.msize {
		//add the curr trans, and generate a microBlock
		var id util.Identifier
		allTxn := append(am.makeTxnSlice(), txn)
		am.currSize = 0
		return true, microblock.NewMicroblock(id, allTxn)

	} else {
		//
		am.txnList.PushBack(txn)
		am.currSize = totalSize
		return false, nil
	}
}

// AddMicroblock adds a microblock into a FIFO queue
// return an err if the queue is full (memsize)
func (am *MemPool) AddMicroblock(mb *microblock.MicroBlock) error {
	am.mu.Lock()
	defer am.mu.Unlock()
	//if am.microblocks.Len() >= am.memsize {
	//	return errors.New("the memory queue is full")
	//}
	_, exists := am.microblockMap[mb.Hash]
	if exists {
		return nil
	}
	pm := &microblock.PendingMicroblock{
		Microblock: mb,
		AckMap:     make(map[util.NodeID]struct{}),
	}
	pm.AckMap[mb.Sender] = struct{}{}
	am.microblockMap[mb.Hash] = mb

	//check if there are some acks of this microblock arrived before
	buffer, received := am.ackBuffer[mb.Hash]
	if received {
		// if so, add these ack to the pendingblocks
		for id, _ := range buffer {
			//am.pendingMicroblocks[mb.Hash].ackMap[ack] = struct{}{}
			pm.AckMap[id] = struct{}{}
		}
		if len(pm.AckMap) >= am.threshhold {
			if _, exists = am.stableMBs[mb.Hash]; !exists {
				am.stableMicroblocks.PushBack(mb)
				am.stableMBs[mb.Hash] = struct{}{}
				delete(am.pendingMicroblocks, mb.Hash)
				//log.Debugf("microblock id: %x becomes stable from buffer", mb.Hash)
			}
		} else {
			am.pendingMicroblocks[mb.Hash] = pm
		}
	} else {
		am.pendingMicroblocks[mb.Hash] = pm
	}
	return nil
}

// AddAck adds an ack and push a microblock into the stableMicroblocks queue if it receives enough acks
func (am *MemPool) AddAck(ack *microblock.Ack) {
	am.mu.Lock()
	defer am.mu.Unlock()
	target, received := am.pendingMicroblocks[ack.MicroblockID]
	//check if the ack arrives before the microblock
	if received {
		target.AckMap[ack.Receiver] = struct{}{}
		if len(target.AckMap) >= am.threshhold {
			// logger.Info().Msgf("One Microblock has received enough acks, current holds %d on Mb %x , ready to propose",len(target.AckMap),target.Microblock.Hash)
			if _, exists := am.stableMBs[target.Microblock.Hash]; !exists {
				am.stableMicroblocks.PushBack(target.Microblock)
				am.stableMBs[target.Microblock.Hash] = struct{}{}
				delete(am.pendingMicroblocks, ack.MicroblockID)
			}
		}
	} else {
		//ack arrives before microblock, record the number of ack received before microblock
		//let the addMicroblock do the rest.
		_, exist := am.ackBuffer[ack.MicroblockID]
		if exist {
			am.ackBuffer[ack.MicroblockID][ack.Receiver] = ack.Signature
		} else {
			temp := make(map[util.NodeID]util.Signature, 0)
			temp[ack.Receiver] = ack.Signature
			am.ackBuffer[ack.MicroblockID] = temp
		}
	}
}

// GeneratePayload generates a list of microblocks according to bsize
// if the remaining microblocks is less than bsize then return all
func (am *MemPool) GeneratePayload() *microblock.Payload {
	var batchSize int
	am.mu.Lock()
	defer am.mu.Unlock()
	sigMap := make(map[util.Identifier]map[util.NodeID]util.Signature, 0)

	if am.stableMicroblocks.Len() >= am.bsize {
		batchSize = am.bsize
	} else {
		batchSize = am.stableMicroblocks.Len()
	}
	microblockList := make([]*microblock.MicroBlock, 0)

	for i := 0; i < batchSize; i++ {
		mb := am.front()
		if mb == nil {
			break
		}
		//log.Debugf("microblock id: %x is deleted from mempool when proposing", mb.Hash)
		microblockList = append(microblockList, mb)

		sigs := make(map[util.NodeID]util.Signature, 0)
		count := 0
		for id, sig := range am.ackBuffer[mb.Hash] {
			count++
			sigs[id] = sig
			// @TODO
			// if count == config.Configuration.Q {
			// 	break
			// }
		}
		sigMap[mb.Hash] = sigs
	}

	return microblock.NewPayload(microblockList, sigMap) // payload 带有mb以及他们相关的ack收集信息以便向其他人证明信息质量
}

// CheckExistence checks if the referred microblocks in the proposal exists
// in the mempool and return missing ones if there's any
// return true if there's no missing transactions

// RemoveMicroblock removes reffered microblocks from the mempool
func (am *MemPool) RemoveMicroblock(id util.Identifier) error {
	am.mu.Lock()
	defer am.mu.Unlock()
	_, exists := am.microblockMap[id]
	if exists {
		delete(am.microblockMap, id)
	}
	_, exists = am.stableMBs[id]
	if exists {
		delete(am.stableMBs, id)
	}
	return nil
}

// FindMicroblock finds a reffered microblock
func (am *MemPool) FindMicroblock(id util.Identifier) (bool, *microblock.MicroBlock) {
	am.mu.Lock()
	defer am.mu.Unlock()
	mb, found := am.microblockMap[id]
	return found, mb
}

func (am *MemPool) CheckExistence(p *microblock.Proposal) (bool, []util.Identifier) {
	id := make([]util.Identifier, 0)
	return false, id
}

// FillProposal pulls microblocks from the mempool and build a pending block,
// a pending block should include the proposal, micorblocks that already exist,
// and a missing list if there's any
func (am *MemPool) FillProposal(p *microblock.Proposal) *microblock.PendingBlock {
	am.mu.Lock()
	defer am.mu.Unlock()
	existingBlocks := make([]*microblock.MicroBlock, 0)
	missingBlocks := make(map[util.Identifier]struct{}, 0)
	for _, id := range p.HashList {
		found := false
		_, exists := am.pendingMicroblocks[id]
		if exists {
			found = true
			existingBlocks = append(existingBlocks, am.pendingMicroblocks[id].Microblock)
			delete(am.pendingMicroblocks, id)
			//log.Debugf("microblock id: %x is deleted from pending when filling", id)
		}
		for e := am.stableMicroblocks.Front(); e != nil; e = e.Next() {
			// do something with e.Value
			mb := e.Value.(*microblock.MicroBlock)
			if mb.Hash == id {
				existingBlocks = append(existingBlocks, mb)
				found = true
				am.stableMicroblocks.Remove(e)
				//log.Debugf("microblock id: %x is deleted from stable when filling", mb.Hash)
				break
			}
		}
		if !found {
			missingBlocks[id] = struct{}{}
		}
	}
	return microblock.NewPendingBlock(p, missingBlocks, existingBlocks)
}

func (am *MemPool) IsStable(id util.Identifier) bool {
	am.mu.Lock()
	defer am.mu.Unlock()
	_, exists := am.stableMBs[id]
	if exists {
		return true
	}
	return false
}

func (am *MemPool) TotalTx() int64 {
	return am.totalTx
}

func (am *MemPool) RemainingTx() int64 {
	return int64(am.txnList.Len())
}

func (am *MemPool) TotalMB() int64 {
	am.mu.Lock()
	defer am.mu.Unlock()
	return int64(len(am.microblockMap))
}

func (am *MemPool) RemainingMB() int64 {
	am.mu.Lock()
	defer am.mu.Unlock()
	return int64(len(am.pendingMicroblocks) + am.stableMicroblocks.Len())
}

func (am *MemPool) AckList(id util.Identifier) []util.NodeID {
	am.mu.Lock()
	defer am.mu.Unlock()
	pmb, exists := am.pendingMicroblocks[id]
	if exists {
		nodes := make([]util.NodeID, 0, len(pmb.AckMap))
		for k, _ := range pmb.AckMap {
			nodes = append(nodes, k)
		}
		return nodes
	}
	return nil
}

func (am *MemPool) front() *microblock.MicroBlock {
	if am.stableMicroblocks.Len() == 0 {
		return nil
	}
	ele := am.stableMicroblocks.Front()
	val, ok := ele.Value.(*microblock.MicroBlock)
	if !ok {
		return nil
	}
	am.stableMicroblocks.Remove(ele)
	return val
}

func (am *MemPool) makeTxnSlice() []*request.Request {
	allTxn := make([]*request.Request, 0)
	for am.txnList.Len() > 0 {
		e := am.txnList.Front()
		allTxn = append(allTxn, e.Value.(*request.Request))
		am.txnList.Remove(e)
	}
	return allTxn
}
