// Copyright 2022 IBM Corp. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package request

import (
	"fmt"
	"sync"

	// "sync"
	// "sync/atomic"

	// "github.com/hyperledger-labs/mirbft/config"
	"github.com/hyperledger-labs/mirbft/crypto"
	"github.com/hyperledger-labs/mirbft/membership"
	"github.com/hyperledger-labs/mirbft/messenger"

	// "github.com/hyperledger-labs/mirbft/membership"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
	logger "github.com/rs/zerolog/log"
)

// Represents a batch of requests.
type Batch struct {
	Sn         int
	MBHashList [][]byte
	SigMap     map[util.Identifier]map[int32][]byte
	BucketId   int
}

type FilledBatch struct {
	Requests []*Request
}

// fillBatch 会对Batch进行填充，调用者需要自己开启线程保证不要阻塞
func (b *Batch) FillBatch(proposerId int) *FilledBatch {
	var wg sync.WaitGroup
	newFilledBatch := &FilledBatch{
		Requests: make([]*Request, 0),
	}
	// 初始化 BucketMBHashList 的 map
	bucketID := b.BucketId
	
	defer wg.Done()
	var MissingMBList []*pb.Identifier
	pendingBlock := Buckets[bucketID].Mempool.FillProposal(b.MBHashList)
	block := pendingBlock.CompleteBlock()
	if block != nil {
		logger.Debug().Msgf("[%v] a block is ready, id: %x", membership.OwnID, b.Sn)
		for _, mb :=  range block.Payload.MicroblockList{
			newFilledBatch.Requests = append(newFilledBatch.Requests, mb.Txns...) 
		}
	}
	PendingBlockMap[b.Sn] = pendingBlock
	logger.Debug().Msgf("%v microblocks are missing in id: %x", len(pendingBlock.MissingMap), b.Sn)
	for mbid, _ := range pendingBlock.MissingMap {
		MissingMBs[mbid] = b.Sn
		MissingMBList = append(MissingMBList, ToProtoIdentifier(mbid))
		logger.Debug().Msgf("[%v] a microblock is missing, id: %x", membership.OwnID, mbid)
	}
	wg.Add(1)
	// 发送Req请求召回missing entry
	missingRequest := pb.MissingMBRequest{
		RequesterId:   membership.OwnID,
		BucketId:      int32(b.BucketId),
		Sn:            int32(b.Sn),
		MissingMbList: MissingMBList,
	}
	msg := &pb.ProtocolMessage{
		SenderId: int32(membership.OwnID),
		Msg: &pb.ProtocolMessage_MissingMicroblockRequest{
			MissingMicroblockRequest: &missingRequest,
		},
	}
	messenger.EnqueueMsg(msg, int32(proposerId))
	return nil
}

// ATTENTION: access to InFLight field is not atomic.
// TODO: do we need to make it atomic? currently the orderer processes all messages pertaining an instance sequentially.
// TODO: is it possible to concurrently access the same request from different instances?

// Marks all requests in the batch as "in flight"
// func (b *Batch) MarkInFlight() {
// 	for _, req := range b.Requests {
// 		req.InFlight = true
// 	}
// }

// Checks if the batch contains "in flight" requests.
// If the batch has "in flight" requests the method returns an error.
// 检测是否有Inflight的请求
// func (b *Batch) CheckInFlight() error {
// 	for _, req := range b.Requests {
// 		if req.InFlight {
// 			return fmt.Errorf("request %d from %d is in flight", req.Msg.RequestId.ClientSn, req.Msg.RequestId.ClientId)
// 		}
// 	}
// 	return nil
// }

// Checks if the requests in the batch match a specific bucket.
// If there exists some request that does not match the bucket id the method returns an error.
// 检查Batch中的请求是否都有激活的桶作为归属
func (b *Batch) CheckBucket(activeBuckets []int) error {
	bucketID := b.BucketId
	found := false
	for _, b := range activeBuckets {
		if b == bucketID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Batch bucket  does not match any active bucket, should be in bucket %d", bucketID)
	}
	return nil
}

// 检查来自CLient的签名
func (b *Batch) CheckSignatures() error {
	if batchVerifierFunc(b) {
		return nil

	} else {
		return fmt.Errorf("batch signature verification failed")
	}
}

// Creates a batch from a protobuf message and tries to add the requests of the message to their buffer.
// If the request is not added successfully (Add returns nil) this method also returns nil
// 利用负载装填一个batch结构并检查签名
// 改动 我们在用pb.Batch的信息还原一个Batch 对于我们装填了MB的Batch来说，似乎验证MB本身 没有意义，因为还是MB的安全是由PAB保证的
func NewBatch(msg *pb.Batch) *Batch {

	logger.Debug().Int("nMB", len(msg.MbHashList)).Msg("Creating new Batch.")

	newBatch := &Batch{
		MBHashList: make([][]byte, len(msg.MbHashList), len(msg.MbHashList)),
		SigMap:     make(map[util.Identifier]map[int32][]byte, len(msg.SigMap)),
		BucketId:   -1,
	}

	newBatch.MBHashList = msg.MbHashList
	for key, value := range msg.SigMap {
		newBatch.SigMap[util.BytesToIdentifier([]byte(key))] = value.MicroblockSigmap
	}
	newBatch.BucketId = int(msg.BucketId)

	// Check signatures of the requests in the new batch.
	// if config.Config.SignRequests {
	// 	if err := newBatch.CheckSignatures(); err != nil {
	// 		logger.Fatal().Err(err).Msg("Invalid signature in new batch.")
	// 		// TODO: Instead of crashing, just return nil.
	// 		return nil
	// 	}
	// }

	return newBatch
}



// 利用负载装填一个batch结构并检查签名
// 改动 我们在用pb.Batch的信息还原一个Batch 对于我们装填了MB的Batch来说，似乎验证MB本身 没有意义，因为还是MB的安全是由PAB保证的
func NewFilledBatch(msg *pb.FilledBatch) *FilledBatch {

	logger.Debug().Int("nMB", len(msg.Requests)).Msg("Creating new Batch.")

	newBatch := &FilledBatch{
		Requests: make([]*Request, 0),
	}
	for _, reqMsg := range msg.Requests{
		req := &Request{
			Msg:      reqMsg,
			Digest:   Digest(reqMsg),
			Buffer:   getBuffer(reqMsg.RequestId.ClientId),
			Bucket:   getBucket(reqMsg),
			Verified: false, // signature has not yet been verified
			InFlight: false, // request has not yet been proposed (an identical one might have been, though, in which case we discard this request object)
			Next:     nil,   // This request object is not part of a bucket list.
			Prev:     nil,
		}
		newBatch.Requests = append(newBatch.Requests, req)
	}

	// Check signatures of the requests in the new batch.
	// if config.Config.SignRequests {
	// 	if err := newBatch.CheckSignatures(); err != nil {
	// 		logger.Fatal().Err(err).Msg("Invalid signature in new batch.")
	// 		// TODO: Instead of crashing, just return nil.
	// 		return nil
	// 	}
	// }

	return newBatch
}

func (b *FilledBatch) Message() *pb.FilledBatch{
	newFilledBatchMsg := &pb.FilledBatch{
		Requests: make([]*pb.ClientRequest, 0),
	}
	for _, req := range b.Requests{
		reqMsg := &pb.ClientRequest{
				RequestId: req.Msg.RequestId,
				Payload: req.Msg.Payload,
				Pubkey:  req.Msg.Pubkey,
				Signature: req.Msg.Signature,
			}
		newFilledBatchMsg.Requests = append(newFilledBatchMsg.Requests, reqMsg)
	}
	return newFilledBatchMsg

}
// Returns a protobuf message containing this Batch.
// 将Batch还原成ReqMsg
func (b *Batch) Message() *pb.Batch {
	// Create empty Batch message
	msg := pb.Batch{
		MbHashList: make([][]byte, len(b.MBHashList), len(b.MBHashList)),
		SigMap:     make(map[string]*pb.MBSig),
		BucketId:   -1,
	}

	// Populate Batch message with request messages
	for i, MBHash := range b.MBHashList {
		msg.MbHashList[i] = MBHash
	}
	msg.BucketId = int32(b.BucketId)
	for key, value := range b.SigMap {
		pbValue := pb.MBSig{
			MicroblockSigmap: value,
		}
		msg.SigMap[string(util.IdentifierToBytes(key))] = &pbValue
	}

	// Return final Batch message
	return &msg
}

// Returns requests in the batch in their buckets after an unsuccessful proposal.
// TODO: Optimization: First group the requests by bucket and then prepend each group at once.
// 将未能Inflight的Req重新取回，应该是在提议后立即调用检查
// @TODO 在MB下，我们还需要管Resurrect，如何将共识失败的MB还原并等待重新提议？ 总之这不是原则问题，我们姑且忽略
func (b *Batch) Resurrect() {
	// for _, req := range b.Requests {
	// 	req.InFlight = false
	// 	req.Bucket.Prepend(req)
	// }
}

// 签名确认应当在收到MB 时回复ACK前完成 ，所以以下三种验证方式均不再需要
// @TODO 在ACK前完成签名验证
func checkSignaturesSequential(b *Batch) bool {
	// for _, req := range b.Requests {
	// 	if !req.Verified {
	// 		if err := crypto.CheckSig(req.Digest,
	// 			membership.ClientPubKey(req.Msg.RequestId.ClientId),
	// 			req.Msg.Signature); err != nil {
	// 			logger.Warn().
	// 				Err(err).
	// 				Int32("clSn", req.Msg.RequestId.ClientSn).
	// 				Int32("clId", req.Msg.RequestId.ClientId).
	// 				Msg("Invalid request signature.")

	// 			return false
	// 		} else {
	// 			req.Verified = true
	// 		}
	// 	}
	// }

	return true
}

func checkSignaturesParallel(b *Batch) bool {
	// var wg sync.WaitGroup
	// wg.Add(len(b.Requests))
	// invalidReqs := int32(0)

	// for _, r := range b.Requests {
	// 	if !r.Verified {
	// 		go func(req *Request) {
	// 			if err := crypto.CheckSig(req.Digest,
	// 				membership.ClientPubKey(req.Msg.RequestId.ClientId),
	// 				req.Msg.Signature); err != nil {
	// 				logger.Warn().
	// 					Err(err).
	// 					Int32("clSn", req.Msg.RequestId.ClientSn).
	// 					Int32("clId", req.Msg.RequestId.ClientId).
	// 					Msg("Invalid request signature.")
	// 				atomic.AddInt32(&invalidReqs, 1)
	// 			} else {
	// 				req.Verified = true
	// 			}
	// 			wg.Done()
	// 		}(r)
	// 	} else {
	// 		wg.Done()
	// 	}
	// }
	// wg.Wait()
	return true
}

func checkSignaturesExternal(b *Batch) bool {
	// verifiedChan := make(chan *Request, len(b.Requests))
	// invalidReqs := 0

	// // Write all the unverified requests in the verifier channel.
	// verifying := 0
	// for _, r := range b.Requests {
	// 	if !r.Verified {
	// 		verifying++
	// 		r.VerifiedChan = verifiedChan
	// 		verifierChan <- r
	// 	}
	// }

	// // Wait until the verifiers process all the requests and write them in the verifiedChan
	// for verifying > 0 {
	// 	verifying--
	// 	req := <-verifiedChan
	// 	req.VerifiedChan = nil
	// 	if !req.Verified {
	// 		logger.Warn().
	// 			Int32("clSn", req.Msg.RequestId.ClientSn).
	// 			Int32("clId", req.Msg.RequestId.ClientId).
	// 			Msg("Request signature verification failed.")
	// 		invalidReqs++
	// 	}
	// }

	// return invalidReqs == 0
	return true
}

// 返回Batch哈希
func BatchDigest(batch *pb.Batch) []byte {
	metadata := make([]byte, 0, 0)
	MBDigests := make([][]byte, len(batch.MbHashList), len(batch.MbHashList))
	for i, MBHash := range batch.MbHashList {
		MBDigests[i] = MBHash
	}
	return crypto.ParallelDataArrayHash(append(MBDigests, crypto.Hash(metadata)))
}
