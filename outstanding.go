/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	"container/list"
	"fmt"

	pb "github.com/IBM/mirbft/mirbftpb"
)

func newOutstandingReqs(clientWindows *clientWindows, networkState *pb.NetworkState) *allOutstandingReqs {
	ao := &allOutstandingReqs{
		numBuckets:          uint64(networkState.Config.NumberOfBuckets),
		buckets:             map[BucketID]*bucketOutstandingReqs{},
		correctRequests:     map[string]*pb.ForwardRequest{},
		outstandingRequests: map[string]*sequence{},
		clientWindows:       clientWindows,
	}

	numBuckets := int(networkState.Config.NumberOfBuckets)

	for i := BucketID(0); i < BucketID(numBuckets); i++ {
		bo := &bucketOutstandingReqs{
			clients: map[uint64]*clientOutstandingReqs{},
		}
		ao.buckets[i] = bo

		for _, client := range networkState.Clients {
			var skipRequests map[uint64]struct{}
			bm := bitmask(client.CommittedMask)
			var firstUncommitted uint64
			var firstOffset int
			for j := 0; j < numBuckets; j++ {
				reqNo := client.LowWatermark + uint64(j)
				if clientReqToBucket(client.Id, reqNo, networkState.Config) == i {
					firstUncommitted = reqNo
					firstOffset = j
					break
				}
			}

			for bitIndex := firstOffset; bitIndex < bm.bits(); bitIndex += numBuckets {
				if !bm.isBitSet(bitIndex) {
					continue
				}

				reqNo := client.LowWatermark + uint64(bitIndex)
				if reqNo == firstUncommitted {
					firstUncommitted += uint64(numBuckets)
					continue
				}

				if skipRequests == nil {
					skipRequests = map[uint64]struct{}{}
				}

				skipRequests[reqNo] = struct{}{}
			}

			bo.clients[client.Id] = &clientOutstandingReqs{
				nextReqNo:    firstUncommitted,
				skipRequests: skipRequests,
			}
		}
	}

	ao.advanceRequests() // Note, this can return no actions as no sequences have allocated

	return ao
}

type allOutstandingReqs struct {
	numBuckets          uint64
	buckets             map[BucketID]*bucketOutstandingReqs
	clientWindows       *clientWindows
	lastCorrectReq      *list.Element
	correctRequests     map[string]*pb.ForwardRequest
	outstandingRequests map[string]*sequence
}

type bucketOutstandingReqs struct {
	clients map[uint64]*clientOutstandingReqs // TODO, obvious optimization is to make this active clients and initialize this lazily
}

type clientOutstandingReqs struct {
	nextReqNo    uint64
	skipRequests map[uint64]struct{}
}

func (ao *allOutstandingReqs) advanceRequests() *Actions {
	actions := &Actions{}
	for {
		var nextCorrectReq *list.Element
		if ao.lastCorrectReq == nil {
			nextCorrectReq = ao.clientWindows.correctList.Front()
		} else {
			nextCorrectReq = ao.lastCorrectReq.Next()
		}

		if nextCorrectReq == nil {
			return actions
		}

		ao.lastCorrectReq = nextCorrectReq

		fr := nextCorrectReq.Value.(*pb.ForwardRequest)
		key := string(fr.Digest)

		if seq, ok := ao.outstandingRequests[key]; ok {
			delete(ao.outstandingRequests, key)
			actions.concat(seq.satisfyOutstanding(fr))
			continue
		}

		ao.correctRequests[key] = fr
	}
}

func (ao *allOutstandingReqs) applyBatch(bucket BucketID, batch []*pb.ForwardRequest) error {
	bo, ok := ao.buckets[bucket]
	if !ok {
		panic("dev sanity test")
	}

	for _, req := range batch {
		co, ok := bo.clients[req.Request.ClientId]
		if !ok {
			return fmt.Errorf("no such client")
		}

		if co.nextReqNo != req.Request.ReqNo {
			return fmt.Errorf("expected ClientId=%d next request for Bucket=%d to have ReqNo=%d but got ReqNo=%d", req.Request.ClientId, bucket, co.nextReqNo, req.Request.ReqNo)
		}

		co.nextReqNo += ao.numBuckets
	}

	return nil
}

// TODO, bucket probably can/should be stored in the *sequence
func (ao *allOutstandingReqs) applyAcks(bucket BucketID, seq *sequence, batch []*pb.RequestAck) (*Actions, error) {
	bo, ok := ao.buckets[bucket]
	if !ok {
		panic("dev sanity test")
	}

	outstandingReqs := map[string]int{}
	forwardReqs := make([]*pb.ForwardRequest, len(batch))

	for i, req := range batch {
		co, ok := bo.clients[req.ClientId]
		if !ok {
			return nil, fmt.Errorf("no such client")
		}

		if co.nextReqNo != req.ReqNo {
			return nil, fmt.Errorf("expected ClientId=%d next request for Bucket=%d to have ReqNo=%d but got ReqNo=%d", req.ClientId, bucket, co.nextReqNo, req.ReqNo)
		}

		key := string(req.Digest)
		if fr, ok := ao.correctRequests[key]; ok {
			delete(ao.correctRequests, key)
			forwardReqs[i] = fr
		} else {
			ao.outstandingRequests[key] = seq
			outstandingReqs[key] = i
		}

		co.nextReqNo += ao.numBuckets
	}

	return seq.allocate(batch, forwardReqs, outstandingReqs), nil
}
