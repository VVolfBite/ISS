package mempool

import (
	"container/list"
	"fmt"
	"sync"
	"log"
	
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
)

type MicroBlock struct{
	mbId *pb.RequestID
}

type MemPool struct {
	microblocks   *list.List
	txnList       *list.List
	microblockMap map[*pb.RequestID]*MicroBlock
	bsize         int // number of microblocks in a proposal
	msize         int // byte size of transactions in a microblock
	memsize       int // number of microblocks in mempool
	currSize      int
	totalTx       int
	mu            sync.Mutex
}

// NewMemPool creates a new naive mempool
func NewMemPool(bsize int,msize int, memsize int) *MemPool {
	return &MemPool{
		bsize:         bsize,
		msize:         msize,
		memsize:       memsize,
		microblocks:   list.New(),
		microblockMap: map[*pb.RequestID]*MicroBlock{},
		currSize:      0,
		txnList:       list.New(),
	}
}

// AddTxn adds a transaction and returns a microblock if msize is reached
// then the contained transactions should be deleted
func (nm *MemPool) AddTxn(txn *pb.ClientRequest) (bool, *MicroBlock) {
	// mempool is full
	if nm.RemainingTx() >= int(nm.memsize) {
		log.Printf("mempool's tx list is full")
		return false, nil
	}
	if nm.RemainingMB() >= int(nm.memsize) {
		log.Printf("mempool's mb list is full")
		return false, nil
	}

	// get the size of the structure. txn is the pointer.
	tranSize := util.SizeOf(txn)
	totalSize := tranSize + nm.currSize

	if tranSize > nm.msize {
		log.Printf("No memory to hold the txn, as txnsize :%d and msize:%d",tranSize,nm.msize)
		return false, nil
	}
	nm.totalTx++
	if totalSize > nm.msize {
		//do not add the curr trans, and generate a microBlock
		//set the currSize to curr trans, since it is the only one does not add to the microblock
		var id *pb.RequestID = txn.RequestId
		nm.currSize = tranSize
		// newBlock := blockchain.NewMicroblock(id, nm.makeTxnSlice())
		newBlock := &MicroBlock{mbId: id}
		nm.txnList.PushBack(txn)
		return true, newBlock

	} else if totalSize == nm.msize {
		//add the curr trans, and generate a microBlock
		var id *pb.RequestID = txn.RequestId
		// allTxn := append(nm.makeTxnSlice(), txn)
		nm.currSize = 0
		return true, &MicroBlock{mbId: id}

	} else {
		nm.txnList.PushBack(txn)
		nm.currSize = totalSize
		return false, nil
	}
}

// AddMicroblock adds a microblock into a FIFO queue
// return an err if the queue is full (memsize)
func (nm *MemPool) AddMicroblock(mb *MicroBlock) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	//if nm.microblocks.Len() >= nm.memsize {
	//	return errors.New("the memory queue is full")
	//}
	_, exists := nm.microblockMap[mb.mbId]
	if exists {
		return nil
	}
	nm.microblocks.PushBack(mb)
	nm.microblockMap[mb.mbId] = mb
	return nil
}

// GeneratePayload generates a list of microblocks according to bsize
// if the remaining microblocks is less than bsize then return all
// func (nm *MemPool) GeneratePayload() *blockchain.Payload {
// 	var batchSize int
// 	nm.mu.Lock()
// 	defer nm.mu.Unlock()

// 	if nm.microblocks.Len() >= nm.bsize {
// 		batchSize = nm.bsize
// 	} else {
// 		batchSize = nm.microblocks.Len()
// 	}
// 	microblockList := make([]*MicroBlock, 0)

// 	for i := 0; i < batchSize; i++ {
// 		mb := nm.front()
// 		if mb == nil {
// 			break
// 		}
// 		microblockList = append(microblockList, mb)
// 	}

// 	return blockchain.NewPayload(microblockList, nil)
// }

// CheckExistence checks if the referred microblocks in the proposal exists
// in the mempool and return missing ones if there's any
// return true if there's no missing transactions
// func (nm *MemPool) CheckExistence(p *blockchain.Proposal) (bool, []crypto.Identifier) {
// 	id := make([]crypto.Identifier, 0)
// 	return false, id
// }

// RemoveMicroblock removes reffered microblocks from the mempool
func (nm *MemPool) RemoveMicroblock(id *pb.RequestID) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	_, exists := nm.microblockMap[id]
	if exists {
		delete(nm.microblockMap, id)
		return nil
	}
	return fmt.Errorf("the microblock does not exist, id: %x", id)
}

// FindMicroblock finds a reffered microblock
func (nm *MemPool) FindMicroblock(id *pb.RequestID) (bool, *MicroBlock) {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	mb, found := nm.microblockMap[id]
	return found, mb
}

// FillProposal pulls microblocks from the mempool and build a pending block,
// a pending block should include the proposal, micorblocks that already exist,
// and a missing list if there's any
// func (nm *MemPool) FillProposal(p *blockchain.Proposal) *blockchain.PendingBlock {
// 	nm.mu.Lock()
// 	defer nm.mu.Unlock()
// 	existingBlocks := make([]*blockchain.MicroBlock, 0)
// 	missingBlocks := make(map[crypto.Identifier]struct{}, 0)
// 	for _, id := range p.HashList {
// 		block, found := nm.microblockMap[id]
// 		if found {
// 			existingBlocks = append(existingBlocks, block)
// 			for e := nm.microblocks.Front(); e != nil; e = e.Next() {
// 				// do something with e.Value
// 				mb := e.Value.(*blockchain.MicroBlock)
// 				if mb == block {
// 					nm.microblocks.Remove(e)
// 					break
// 				}
// 			}
// 		} else {
// 			missingBlocks[id] = struct{}{}
// 		}
// 	}
// 	return blockchain.NewPendingBlock(p, missingBlocks, existingBlocks)
// }

func (nm *MemPool) TotalTx() int {
	return nm.totalTx
}

func (nm *MemPool) RemainingTx() int {
	return int(nm.txnList.Len())
}

func (nm *MemPool) TotalMB() int {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	return int(len(nm.microblockMap))
}

func (nm *MemPool) RemainingMB() int {
	nm.mu.Lock()
	defer nm.mu.Unlock()
	return int(nm.microblocks.Len())
}

// func (nm *MemPool) AddAck(ack *blockchain.Ack) {
// }

// func (nm *MemPool) AckList(id crypto.Identifier) []identity.NodeID {
// 	return nil
// }

func (nm *MemPool) IsStable(id *pb.RequestID) bool {
	return false
}

func (nm *MemPool) front() *MicroBlock {
	if nm.microblocks.Len() == 0 {
		return nil
	}
	ele := nm.microblocks.Front()
	val, ok := ele.Value.(*MicroBlock)
	if !ok {
		return nil
	}
	nm.microblocks.Remove(ele)
	return val
}

// func (nm *MemPool) makeTxnSlice() []*message.Transaction {
// 	allTxn := make([]*message.Transaction, 0)
// 	for nm.txnList.Len() > 0 {
// 		e := nm.txnList.Front()
// 		allTxn = append(allTxn, e.Value.(*message.Transaction))
// 		nm.txnList.Remove(e)
// 	}
// 	return allTxn
// }
