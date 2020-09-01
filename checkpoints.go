/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	"bytes"
	"container/list"
	"fmt"
	"sort"

	pb "github.com/IBM/mirbft/mirbftpb"

	"google.golang.org/protobuf/proto"
)

type checkpointState int

const (
	cpsIdle checkpointState = iota
	cpsGarbageCollectable
	cpsPendingReconfig // TODO, implement
	cpsStateTransfer   // TODO, implement
)

type checkpointTracker struct {
	state checkpointState

	highestCheckpoints map[NodeID]uint64
	checkpointMap      map[uint64]*checkpoint
	activeCheckpoints  *list.List

	networkConfig *pb.NetworkState_Config
	persisted     *persisted
	myConfig      *pb.StateEvent_InitialParameters
}

func newCheckpointTracker(persisted *persisted, myConfig *pb.StateEvent_InitialParameters) *checkpointTracker {
	ct := &checkpointTracker{
		highestCheckpoints: map[NodeID]uint64{},
		checkpointMap:      map[uint64]*checkpoint{},
		myConfig:           myConfig,
		persisted:          persisted,
		state:              cpsIdle,
		activeCheckpoints:  list.New(),
	}

	for head := persisted.logHead; head != nil; head = head.next {
		switch d := head.entry.Type.(type) {
		case *pb.Persistent_CEntry:
			cEntry := d.CEntry
			if ct.networkConfig == nil {
				ct.networkConfig = cEntry.NetworkState.Config
			}
			if !proto.Equal(cEntry.NetworkState.Config, ct.networkConfig) {
				// TODO, implement reconfig
				ct.state = cpsPendingReconfig
				panic("reconfig not yet supported")
			}
			cp := ct.checkpoint(cEntry.SeqNo)
			cp.nextState = cEntry.NetworkState
			cp.applyCheckpointMsg(NodeID(myConfig.Id), cEntry.CheckpointValue)
			ct.activeCheckpoints.PushBack(cp)
		}
	}

	if ct.activeCheckpoints.Len() == 0 {
		panic("no checkpoints in log")
	}

	ct.activeCheckpoints.Front().Value.(*checkpoint).stable = true

	ct.garbageCollect()

	return ct
}

func (ct *checkpointTracker) filter(msg *pb.Msg) applyable {
	switch {
	case msg.Type.(*pb.Msg_Checkpoint).Checkpoint.SeqNo < ct.activeCheckpoints.Front().Value.(*checkpoint).seqNo:
		return past
	default:
		return current
		// TODO, have notion of future... but also process
	}
}

func (ct *checkpointTracker) step(source NodeID, msg *pb.Msg) {
	if ct.filter(msg) != current {
		return
	}

	ct.applyMsg(source, msg)
}

func (ct *checkpointTracker) applyMsg(source NodeID, msg *pb.Msg) {
	switch innerMsg := msg.Type.(type) {
	case *pb.Msg_Checkpoint:
		msg := innerMsg.Checkpoint
		ct.applyCheckpointMsg(source, msg.SeqNo, msg.Value)
	default:
		panic(fmt.Sprintf("unexpected bad checkpoint message type %T, this indicates a bug", msg.Type))
	}
}

func (ct *checkpointTracker) garbageCollect() uint64 {
	var highestStable *list.Element
	for el := ct.activeCheckpoints.Front(); el != nil; el = el.Next() {
		cp := el.Value.(*checkpoint)
		if !cp.stable {
			break
		}

		highestStable = el
	}

	for el := highestStable.Prev(); el != nil; el = highestStable.Prev() {
		delete(ct.checkpointMap, el.Value.(*checkpoint).seqNo)
		ct.activeCheckpoints.Remove(el)
	}

	for ct.activeCheckpoints.Len() < 3 {
		nextCpSeq := ct.highWatermark() + uint64(ct.networkConfig.CheckpointInterval)
		ct.activeCheckpoints.PushBack(ct.checkpoint(nextCpSeq))
	}

	ct.state = cpsIdle
	return highestStable.Value.(*checkpoint).seqNo
}

func (ct *checkpointTracker) checkpoint(seqNo uint64) *checkpoint {
	cp, ok := ct.checkpointMap[seqNo]
	if !ok {
		cp = &checkpoint{
			seqNo:           seqNo,
			verifyingConfig: ct.networkConfig,
			persisted:       ct.persisted,
			myConfig:        ct.myConfig,
		}
		ct.checkpointMap[seqNo] = cp
	}

	return cp
}

func (ct *checkpointTracker) highWatermark() uint64 {
	return ct.activeCheckpoints.Back().Value.(*checkpoint).seqNo
}

func (ct *checkpointTracker) lowWatermark() uint64 {
	return ct.activeCheckpoints.Front().Value.(*checkpoint).seqNo
}

func (ct *checkpointTracker) applyCheckpointMsg(source NodeID, seqNo uint64, value []byte) {
	if seqNo < ct.lowWatermark() {
		// We're already past this point
		return
	}

	aboveHighWatermark := seqNo > ct.highWatermark()
	if aboveHighWatermark {
		highest, ok := ct.highestCheckpoints[source]
		if ok && highest <= seqNo {
			return
		}

		ct.highestCheckpoints[source] = seqNo
	}

	cp := ct.checkpoint(seqNo)
	cp.applyCheckpointMsg(source, value)

	if seqNo > ct.lowWatermark() && cp.stable {
		ct.state = cpsGarbageCollectable
		return
	}

	if !aboveHighWatermark {
		return
	}

	// We just added a new entry to our highest checkpoints map,
	// so we need to garbage collect any above window checkpoint
	// references that no node claims is the most current anymore.

	referencedCPs := map[uint64]struct{}{}

	for el := ct.activeCheckpoints.Front(); el != nil; el = el.Next() {
		referencedCPs[el.Value.(*checkpoint).seqNo] = struct{}{}
	}

	for _, seqNo := range ct.highestCheckpoints {
		referencedCPs[seqNo] = struct{}{}
	}

	for seqNo := range ct.checkpointMap {
		if _, ok := referencedCPs[seqNo]; !ok {
			delete(ct.checkpointMap, seqNo)
		}
	}
}

func (ct *checkpointTracker) applyCheckpointResult(seqNo uint64, value []byte, epochConfig *pb.EpochConfig, nextConfig *pb.NetworkState) *Actions {
	return ct.checkpoint(seqNo).applyCheckpointResult(value, epochConfig, nextConfig)
}

func (ct *checkpointTracker) status() []*CheckpointStatus {
	result := make([]*CheckpointStatus, len(ct.checkpointMap))
	i := 0
	for _, cp := range ct.checkpointMap {
		result[i] = cp.status()
		i++
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].SeqNo < result[j].SeqNo
	})

	return result
}

type checkpoint struct {
	seqNo           uint64
	myConfig        *pb.StateEvent_InitialParameters
	verifyingConfig *pb.NetworkState_Config
	persisted       *persisted

	values         map[string][]NodeID
	committedValue []byte
	myValue        []byte
	nextState      *pb.NetworkState
	stable         bool
	obsolete       bool
}

func (cw *checkpoint) applyCheckpointMsg(source NodeID, value []byte) bool {
	if cw.values == nil {
		cw.values = map[string][]NodeID{}
	}

	stateChange := false

	checkpointValueNodes := append(cw.values[string(value)], source)
	cw.values[string(value)] = checkpointValueNodes

	agreements := len(checkpointValueNodes)

	if agreements == someCorrectQuorum(cw.verifyingConfig) {
		cw.committedValue = value
	}

	if source == NodeID(cw.myConfig.Id) {
		cw.myValue = value
	}

	// If I have completed this checkpoint, along with a quorum of the network, and I've not already run this path
	if cw.myValue != nil && cw.committedValue != nil && !cw.stable {
		if !bytes.Equal(value, cw.committedValue) {
			// TODO optionally handle this more gracefully, with state transfer (though this
			// indicates a violation of the byzantine assumptions)
			panic("my checkpoint disagrees with the committed network view of this checkpoint")
		}

		// This checkpoint has enough agreements, including my own, it may now be garbage collectable
		// Note, this must be >= (not ==) because my agreement could come after 2f+1 from the network.
		if agreements >= intersectionQuorum(cw.verifyingConfig) {
			cw.stable = true
			stateChange = true
		}
	}

	if len(checkpointValueNodes) == len(cw.verifyingConfig.Nodes) {
		cw.obsolete = true
		stateChange = true
	}

	return stateChange
}

func (cw *checkpoint) applyCheckpointResult(value []byte, epochConfig *pb.EpochConfig, nextState *pb.NetworkState) *Actions {
	cw.nextState = nextState
	return (&Actions{}).send(
		cw.verifyingConfig.Nodes,
		&pb.Msg{
			Type: &pb.Msg_Checkpoint{
				Checkpoint: &pb.Checkpoint{
					SeqNo: uint64(cw.seqNo),
					Value: value,
				},
			},
		},
	).concat(cw.persisted.addCEntry(&pb.CEntry{
		SeqNo:           cw.seqNo,
		CheckpointValue: value,
		NetworkState:    nextState,
		EpochConfig:     epochConfig,
	}))
}

func (cw *checkpoint) status() *CheckpointStatus {
	maxAgreements := 0
	for _, nodes := range cw.values {
		if len(nodes) > maxAgreements {
			maxAgreements = len(nodes)
		}
	}
	return &CheckpointStatus{
		SeqNo:         cw.seqNo,
		MaxAgreements: maxAgreements,
		NetQuorum:     cw.committedValue != nil,
		LocalDecision: cw.myValue != nil,
	}
}
