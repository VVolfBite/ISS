/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	pb "github.com/IBM/mirbft/mirbftpb"
)

// epochConfig is the information required by the various
// state machines whose state is scoped to an epoch
type epochConfig struct {
	// myConfig is the configuration specific to this node
	myConfig *Config

	// oddities stores counts of suspicious acitivities and logs them.
	oddities *oddities

	// number is the epoch number this config applies to
	number uint64

	// highWatermark is the current maximum seqno that may be processed
	highWatermark SeqNo

	// lowWatermark is the current minimum seqno for which messages are valid
	lowWatermark SeqNo

	// F is the total number of faults tolerated by the network
	f int

	// CheckpointInterval is the number of sequence numbers to commit before broadcasting a checkpoint
	checkpointInterval SeqNo

	// nodes is all the node ids in the network
	nodes []NodeID

	// buckets is a map from bucket ID to leader ID
	buckets map[BucketID]NodeID
}

type epoch struct {
	epochConfig *epochConfig

	buckets map[BucketID]*bucket

	proposer *proposer

	checkpointWindows map[SeqNo]*checkpointWindow

	largestPreprepares map[NodeID]SeqNo
}

func newEpoch(config *epochConfig) *epoch {
	largestPreprepares := map[NodeID]SeqNo{}
	for _, id := range config.nodes {
		largestPreprepares[id] = config.lowWatermark
	}

	buckets := map[BucketID]*bucket{}
	for bucketID := range config.buckets {
		buckets[bucketID] = newBucket(config, bucketID)
	}

	checkpointWindows := map[SeqNo]*checkpointWindow{}
	for seqNo := config.lowWatermark + config.checkpointInterval; seqNo <= config.highWatermark; seqNo += config.checkpointInterval {
		checkpointWindows[seqNo] = newCheckpointWindow(seqNo, config)
	}

	return &epoch{
		epochConfig:        config,
		buckets:            buckets,
		checkpointWindows:  checkpointWindows,
		largestPreprepares: largestPreprepares,
		proposer:           newProposer(config),
	}
}

func (e *epoch) process(preprocessResult PreprocessResult) *Actions {
	bucketID := BucketID(preprocessResult.Cup % uint64(len(e.epochConfig.buckets)))
	nodeID := e.epochConfig.buckets[bucketID]
	if nodeID == NodeID(e.epochConfig.myConfig.ID) {
		return e.proposer.propose(preprocessResult.Proposal.Data)
	}

	if preprocessResult.Proposal.Source == e.epochConfig.myConfig.ID {
		// I originated this proposal, but someone else leads this bucket,
		// forward the message to them
		return &Actions{
			Unicast: []Unicast{
				{
					Target: uint64(nodeID),
					Msg: &pb.Msg{
						Type: &pb.Msg_Forward{
							Forward: &pb.Forward{
								Epoch:  e.epochConfig.number,
								Bucket: uint64(bucketID),
								Data:   preprocessResult.Proposal.Data,
							},
						},
					},
				},
			},
		}
	}

	// Someone forwarded me this proposal, but I'm not responsible for it's bucket
	// TODO, log oddity? Assign it to the wrong bucket? Forward it again?
	return &Actions{}
}

func (e *epoch) Preprepare(source NodeID, seqNo SeqNo, bucket BucketID, batch [][]byte) *Actions {
	actions := &Actions{}

	if e.largestPreprepares[source] < seqNo {
		e.largestPreprepares[source] = seqNo

		farAhead := e.proposer.nextAssigned + e.epochConfig.checkpointInterval

		if seqNo >= farAhead {
			// XXX this is really a kind of heuristic check, to make sure
			// that if the network is advancing without us, possibly because
			// of unbalanced buckets, that we keep up, it's worth formalizing.
			nodesFurtherThanMe := 0
			for _, largest := range e.largestPreprepares {
				// XXX it's possible the node does not lead any bucket, perhaps the map
				// should exclude these?
				if largest >= farAhead {
					nodesFurtherThanMe++
				}
			}

			if nodesFurtherThanMe > e.epochConfig.f {
				actions.Append(e.proposer.noopAdvance())
			}
		}
	}

	actions.Append(e.buckets[bucket].applyPreprepare(seqNo, batch))
	return actions
}

func (e *epoch) Prepare(source NodeID, seqNo SeqNo, bucket BucketID, digest []byte) *Actions {
	return e.buckets[bucket].applyPrepare(source, seqNo, digest)
}

func (e *epoch) Commit(source NodeID, seqNo SeqNo, bucket BucketID, digest []byte) *Actions {
	actions := e.buckets[bucket].applyCommit(source, seqNo, digest)
	if len(actions.Commit) > 0 {
		// XXX this is a moderately hacky way to determine if this commit msg triggered
		// a commit, is there a better way?
		if checkpointWindow, ok := e.checkpointWindows[seqNo]; ok {
			actions.Append(checkpointWindow.Committed(bucket))
		}
	}
	return actions
}

func (e *epoch) moveWatermarks() *Actions {
	for _, bucket := range e.buckets {
		bucket.moveWatermarks()
	}

	return e.proposer.drainQueue()
}

func (e *epoch) checkpointResult(seqNo SeqNo, value, attestation []byte) *Actions {
	checkpointWindow, ok := e.checkpointWindows[seqNo]
	if !ok {
		panic("received an unexpected checkpoint result")
	}
	return checkpointWindow.applyCheckpointResult(value, attestation)
}

func (e *epoch) digest(seqNo SeqNo, bucket BucketID, digest []byte) *Actions {
	return e.buckets[bucket].applyDigestResult(seqNo, digest)
}

func (e *epoch) validate(seqNo SeqNo, bucket BucketID, valid bool) *Actions {
	return e.buckets[bucket].applyValidateResult(seqNo, valid)
}

func (e *epoch) Tick() *Actions {
	return e.proposer.noopAdvance()
}
