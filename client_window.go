/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mirbft

import (
	"container/list"
	"fmt"
	"sort"

	pb "github.com/IBM/mirbft/mirbftpb"
)

type clientWindows struct {
	windows       map[uint64]*clientWindow
	clients       []uint64
	networkConfig *pb.NetworkConfig
	myConfig      *Config
	readyList     *list.List
	readyMap      map[*clientReqNo]*list.Element
}

func newClientWindows(networkConfig *pb.NetworkConfig, myConfig *Config) *clientWindows {
	cws := &clientWindows{
		myConfig:      myConfig,
		networkConfig: networkConfig,
		windows:       map[uint64]*clientWindow{},
		readyList:     list.New(),
		readyMap:      map[*clientReqNo]*list.Element{},
	}

	clientWindowWidth := uint64(100) // XXX this should be configurable

	for _, client := range networkConfig.Clients {
		lowWatermark := client.BucketLowWatermarks[0]
		for _, blw := range client.BucketLowWatermarks {
			if blw < lowWatermark {
				lowWatermark = blw
			}
		}

		clientWindow := newClientWindow(lowWatermark, lowWatermark+clientWindowWidth, networkConfig, myConfig)
		cws.insert(client.Id, clientWindow)
	}

	return cws
}

func (cws *clientWindows) ack(source NodeID, clientID, reqNo uint64, digest []byte) {
	cw, ok := cws.windows[clientID]
	if !ok {
		panic("dev sanity test")
	}

	clientReqNo := cw.ack(source, reqNo, digest)
	if clientReqNo.strongRequest == nil {
		return
	}

	_, ok = cws.readyMap[clientReqNo]
	if ok {
		return
	}

	el := cws.readyList.PushBack(clientReqNo)

	cws.readyMap[clientReqNo] = el
}

func (cws *clientWindows) garbageCollect(seqNo uint64) {
	cwi := cws.iterator()
	for _, cw := cwi.next(); cw != nil; _, cw = cwi.next() {
		cw.garbageCollect(seqNo)
	}

	for el := cws.readyList.Front(); el != nil; el = el.Next() {
		c := el.Value.(*clientReqNo).committed
		if c == nil || *c > seqNo {
			continue
		}

		cws.readyList.Remove(el)
		delete(cws.readyMap, el.Value.(*clientReqNo))
	}
}

func (cws *clientWindows) clientWindow(clientID uint64) (*clientWindow, bool) {
	cw, ok := cws.windows[clientID]
	return cw, ok
}

func (cws *clientWindows) insert(clientID uint64, cw *clientWindow) {
	cws.windows[clientID] = cw
	cws.clients = append(cws.clients, clientID)
	sort.Slice(cws.clients, func(i, j int) bool {
		return cws.clients[i] < cws.clients[j]
	})
}

func (cws *clientWindows) iterator() *clientWindowIterator {
	return &clientWindowIterator{
		clientWindows: cws,
	}
}

type clientWindowIterator struct {
	index         int
	clientWindows *clientWindows
}

func (cwi *clientWindowIterator) next() (uint64, *clientWindow) {
	if cwi.index >= len(cwi.clientWindows.clients) {
		return 0, nil
	}
	client := cwi.clientWindows.clients[cwi.index]
	clientWindow := cwi.clientWindows.windows[client]
	cwi.index++
	return client, clientWindow
}

type clientReqNo struct {
	digests       map[string]*clientRequest
	committed     *uint64
	strongRequest *clientRequest
}

type clientRequest struct {
	digest     []byte
	data       *pb.Request
	agreements map[NodeID]struct{}
}

type clientWindow struct {
	lowWatermark   uint64
	highWatermark  uint64
	clientRequests []*clientReqNo
	clientWaiter   *clientWaiter // Used to throttle clients
	myConfig       *Config
	networkConfig  *pb.NetworkConfig
}

type clientWaiter struct {
	lowWatermark  uint64
	highWatermark uint64
	expired       chan struct{}
}

func newClientWindow(lowWatermark, highWatermark uint64, networkConfig *pb.NetworkConfig, myConfig *Config) *clientWindow {
	return &clientWindow{
		myConfig:       myConfig,
		networkConfig:  networkConfig,
		lowWatermark:   lowWatermark,
		highWatermark:  highWatermark,
		clientRequests: make([]*clientReqNo, int(highWatermark-lowWatermark)+1),
		clientWaiter: &clientWaiter{
			lowWatermark:  lowWatermark,
			highWatermark: highWatermark,
			expired:       make(chan struct{}),
		},
	}
}

func (cw *clientWindow) garbageCollect(maxSeqNo uint64) {
	newRequests := make([]*clientReqNo, int(cw.highWatermark-cw.lowWatermark)+1)
	i := 0
	j := uint64(0)
	copying := false
	for _, request := range cw.clientRequests {
		if request == nil || request.committed == nil || *request.committed > maxSeqNo {
			copying = true
		}

		if copying {
			newRequests[i] = request
			i++
		} else {
			if request.committed == nil {
				panic("this should be initialized if here")
			}
			j++
		}

	}

	cw.lowWatermark += j
	cw.highWatermark += j
	cw.clientRequests = newRequests
	close(cw.clientWaiter.expired)
	cw.clientWaiter = &clientWaiter{
		lowWatermark:  cw.lowWatermark,
		highWatermark: cw.highWatermark,
		expired:       make(chan struct{}),
	}
}

func (cw *clientWindow) ack(source NodeID, reqNo uint64, digest []byte) *clientReqNo {
	if reqNo > cw.highWatermark {
		panic(fmt.Sprintf("unexpected: %d > %d", reqNo, cw.highWatermark))
	}

	if reqNo < cw.lowWatermark {
		panic(fmt.Sprintf("unexpected: %d < %d", reqNo, cw.lowWatermark))
	}

	offset := int(reqNo - cw.lowWatermark)
	crn := cw.clientRequests[offset]
	if crn == nil {
		crn = &clientReqNo{
			digests: map[string]*clientRequest{},
		}
		cw.clientRequests[offset] = crn
	}

	cr, ok := crn.digests[string(digest)]
	if !ok {
		cr = &clientRequest{
			digest:     digest,
			agreements: map[NodeID]struct{}{},
		}
		crn.digests[string(digest)] = cr
	}

	cr.agreements[source] = struct{}{}

	if len(cr.agreements) == intersectionQuorum(cw.networkConfig) {
		crn.strongRequest = cr
	}

	return crn
}

func (cw *clientWindow) allocate(requestData *pb.Request, digest []byte) {
	reqNo := requestData.ReqNo
	if reqNo > cw.highWatermark {
		panic(fmt.Sprintf("unexpected: %d > %d", reqNo, cw.highWatermark))
	}

	if reqNo < cw.lowWatermark {
		panic(fmt.Sprintf("unexpected: %d < %d", reqNo, cw.lowWatermark))
	}

	offset := int(reqNo - cw.lowWatermark)
	if cw.clientRequests[offset] == nil {
		cw.clientRequests[offset] = &clientReqNo{
			digests: map[string]*clientRequest{},
		}
	}

	cr, ok := cw.clientRequests[offset].digests[string(digest)]
	if !ok {
		cr = &clientRequest{
			digest:     digest,
			data:       requestData,
			agreements: map[NodeID]struct{}{},
		}
		cw.clientRequests[offset].digests[string(digest)] = cr
	}

	if len(cr.agreements) == intersectionQuorum(cw.networkConfig) {
		cw.clientRequests[offset].strongRequest = cr
	}
}

func (cw *clientWindow) inWatermarks(reqNo uint64) bool {
	return reqNo <= cw.highWatermark && reqNo >= cw.lowWatermark
}

func (cw *clientWindow) request(reqNo uint64) *clientReqNo {
	if reqNo > cw.highWatermark {
		panic(fmt.Sprintf("unexpected: %d > %d", reqNo, cw.highWatermark))
	}

	if reqNo < cw.lowWatermark {
		panic(fmt.Sprintf("unexpected: %d < %d", reqNo, cw.lowWatermark))
	}

	offset := int(reqNo - cw.lowWatermark)

	return cw.clientRequests[offset]
}

func (cw *clientWindow) status() *ClientWindowStatus {
	allocated := make([]uint64, len(cw.clientRequests))
	for i, request := range cw.clientRequests {
		if request == nil {
			continue
		}
		if request.committed != nil {
			allocated[i] = 2 // TODO, actually report the seqno it committed to
		} else {
			allocated[i] = 1
		}
		// allocated[i] = bytesToUint64(request.preprocessResult.Proposal.Data)
	}

	return &ClientWindowStatus{
		LowWatermark:  cw.lowWatermark,
		HighWatermark: cw.highWatermark,
		Allocated:     allocated,
	}
}
