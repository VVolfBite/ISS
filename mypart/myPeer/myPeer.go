package myPeer

import (
	"context"
	"fmt"
	"github.com/hyperledger-labs/mirbft/mypart/mempool"
	"github.com/hyperledger-labs/mirbft/mypart/microblock"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/request"
	"github.com/hyperledger-labs/mirbft/util"
	logger "github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"time"
	// "google.golang.org/grpc/internal/resolver/passthrough"
)

type MyPeer struct {
	Id              int
	MissingCounts   map[util.NodeID]int
	MissingMBs      map[util.Identifier]util.Identifier          // 缺失的mb，即所有 pending block的快速映射
	ReceivedMBs     map[util.Identifier]struct{}                 // 收到的mb
	PendingBlockMap map[util.Identifier]*microblock.PendingBlock // pending block 是对fill proposal后仍有mb没有获取需要retrive的的block的称呼
	Mempool         *mempool.MemPool
	PeerConnections map[int]*grpc.ClientConn
	IsBusy          bool
}

func NewMyPeer(Id int) (*MyPeer, error) {
	peer := &MyPeer{
		Id:              Id,
		ReceivedMBs:     make(map[util.Identifier]struct{}),
		MissingMBs:      make(map[util.Identifier]util.Identifier),
		MissingCounts:   make(map[util.NodeID]int),
		Mempool:         mempool.NewMemPool(),
		PeerConnections: make(map[int]*grpc.ClientConn),
		IsBusy:          false,
	}

	server := grpc.NewServer()
	pb.RegisterMessengerServer(server, peer)
	address := fmt.Sprintf("localhost:%d", 50000+Id)
	lis, err := net.Listen("tcp", address)

	if err != nil {
		logger.Err(err).Msg(fmt.Sprintf("listen to %d", peer.Id))
		return nil, err
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			logger.Err(err).Msg(fmt.Sprintf("failed to serve: %v", err))
		}
	}()

	return peer, nil
}

// Listen 实现了 Listen RPC 方法
func (p *MyPeer) Listen(srv pb.Messenger_ListenServer) error {
	for {
		msg, err := srv.Recv()
		if err != nil {
			return err
		}
		p.HandleProtocolMsg(msg, srv)

		// 该ack为空值ack
		ack := &pb.BandwidthTestAck{}
		if err := srv.Send(ack); err != nil {
			return err
		}
	}
}

func (p *MyPeer) Request(srv pb.Messenger_RequestServer) error {
	for {
		req, err := srv.Recv()
		if err != nil {
			return err
		}
		p.HandleRequest(req, srv)

		// 创建 ClientResponse 消息
		resp := &pb.ClientResponse{
			ClientSn: req.RequestId.ClientSn,
			OrderSn:  int32(p.Id),
		}

		// 发送 ClientResponse 响应
		if err := srv.Send(resp); err != nil {
			return err
		}
	}
}

// Buckets 实现了 Buckets RPC 方法
func (p *MyPeer) Buckets(srv pb.Messenger_BucketsServer) error {
	// 实现 Buckets RPC 方法的逻辑
	return nil
}

// ConnectToPeer 尝试连接到单个对等节点
func (p *MyPeer) ConnectToPeer(Id int) error {
	address := fmt.Sprintf("localhost:%d", 50000+Id)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	p.PeerConnections[Id] = conn
	return nil
}

// ConnectToAllPeers 尝试连接从 start 到 end 的所有对等节点
func (p *MyPeer) ConnectToAllPeers(start int, end int) error {
	for Id := start; Id <= end; Id++ {

		if err := p.ConnectToPeer(Id); err != nil {
			logger.Err(err).Msgf("Peer %d failed to connect to peer %d", p.Id, Id)
		}
	}
	return nil
}

// Request的收发部分，主要调用对方的Request

func (p *MyPeer) SendRequest(Id int, reqMsg *pb.ClientRequest) error {
	conn, ok := p.PeerConnections[Id]
	if !ok {
		return fmt.Errorf("peer %d is not connected", Id)
	}

	client := pb.NewMessengerClient(conn)
	srv, err := client.Request(context.Background())
	if err != nil {
		return err
	}
	if err := srv.Send(reqMsg); err != nil {
		return err
	}
	return nil
}

func (p *MyPeer) BroadcastRequest(reqMsg *pb.ClientRequest) error {
	for id := range p.PeerConnections {
		err := p.SendRequest(id, reqMsg)
		if err != nil {
			logger.Err(err).Msgf("Error sending request to peer %d", id)
		}
	}
	return nil
}

// HandleRequest 处理请求，简单地调用 Mempool 的 ReceiveTxn 方法
func (p *MyPeer) HandleRequest(reqMsg *pb.ClientRequest, srv pb.Messenger_RequestServer) {
	// 将请求添加到内存池
	var req = &request.Request{
		Msg:      reqMsg,
		Digest:   request.Digest(reqMsg),
		Buffer:   nil,
		Bucket:   nil,
		Verified: false, // signature has not yet been verified
		InFlight: false, // request has not yet been proposed (an identical one might have been, though, in which case we discard this request object)
		Next:     nil,   // This request object is not part of a bucket list.
		Prev:     nil,
	}
	digest := util.BytesToIdentifier(req.Digest)
	if _, exist := p.ReceivedMBs[digest]; exist {
		logger.Info().Msgf("Request with digest %v already exists, skipping...", digest)
		return
	}

	p.ReceivedMBs[digest] = struct{}{}
	isbuilt, mb := p.Mempool.AddTxn(req)
	if isbuilt {
		p.Mempool.AddMicroblock(mb)
		mb.Timestamp = time.Now()

		pMsg := &pb.ProtocolMessage{
			SenderId: int32(p.Id),
			Msg: &pb.ProtocolMessage_Microblock{
				Microblock: microblock.ToProtoMicroBlock(mb),
			},
		}
		// @TODO 负载均衡则可以代理出去
		if p.IsBusy {
			p.SenderLoadBalance(pMsg)
		} else {
			p.BroadcastProtocolMsg(pMsg)
		}

	}
}

// Protocol Message 部分 ，主要调用对方的Listen
func (p *MyPeer) SendProtocolMsg(Id int, protocolMsg *pb.ProtocolMessage) error {

	conn, ok := p.PeerConnections[Id]
	if !ok {
		return fmt.Errorf("peer %d is not connected", Id)
	}

	client := pb.NewMessengerClient(conn)
	srv, err := client.Listen(context.Background())
	if err != nil {
		return err
	}
	if err := srv.Send(protocolMsg); err != nil {
		return err
	}
	return nil
}

func (p *MyPeer) BroadcastProtocolMsg(protocolMsg *pb.ProtocolMessage) error {
	for id := range p.PeerConnections {
		err := p.SendProtocolMsg(id, protocolMsg)
		if err != nil {
			logger.Err(err).Msgf("Error sending request to peer %d", id)
		}
	}
	return nil
}

func (p *MyPeer) HandleProtocolMsg(protocolMsg *pb.ProtocolMessage, srv pb.Messenger_ListenServer) {
	switch msg := protocolMsg.Msg.(type) {

	case *pb.ProtocolMessage_Microblock:
		// 接收到 Microblock 消息
		mb := msg.Microblock
		p.HandleMicroblock(microblock.FromProtoMicroBlock(mb))
		// ack相关内容应当在handler中处理
	case *pb.ProtocolMessage_MicroblockAck:
		// 接收到 Ack 消息
		ack := msg.MicroblockAck
		p.HandleAck(microblock.FromProtoAck(ack))
	case *pb.ProtocolMessage_MissingMicroblockRequest:
		// 接收到
		missmbReq := msg.MissingMicroblockRequest
		p.HandleMissingMicroblockRequest(microblock.FromProtoMissingMBRequest(missmbReq))
	default:
		// 其他类型的消息，暂不处理
		logger.Printf("Received unknown message type: %T\n", msg)
	}
}

// bucket部分还不需要实现

func (p *MyPeer) HandleMicroblock(mb *microblock.MicroBlock) {
	// logger.Info().Msgf("Received microblock message from sender %d with Hash %x", mb.Sender, mb.Hash)
	_, exist := p.ReceivedMBs[mb.Hash]

	if exist {
		return
	}
	p.ReceivedMBs[mb.Hash] = struct{}{}
	mb.FutureTimestamp = time.Now()
	proposalID, exists := p.MissingMBs[mb.Hash]
	if exists {
		pd, exists := p.PendingBlockMap[proposalID]
		if exists {
			block := pd.AddMicroblock(mb)
			if block != nil {
				delete(p.PendingBlockMap, mb.ProposalID)
				delete(p.MissingMBs, mb.Hash)
			}
		}
	} else {
		// 不然对方却的也是自己缺少的，加入mempool
		err := p.Mempool.AddMicroblock(mb)
		if err != nil {
		}
		// ack
		if !mb.IsRequested {
			//if config.Configuration.MemType == "time" {
			//	r.Send(mb.Sender, ack)
			//} else {
			//	r.Broadcast(ack)
			//}
			ack := &microblock.Ack{
				Receiver:     util.NodeID(p.Id),
				MicroblockID: mb.Hash,
				Signature:    nil,
			}
			msg := &pb.ProtocolMessage{
				SenderId: int32(p.Id),
				Msg: &pb.ProtocolMessage_MicroblockAck{
					MicroblockAck: microblock.ToProtoAck(ack),
				},
			}
			if mb.Sender != util.NodeID(p.Id) {
				p.SendProtocolMsg(int(mb.Sender), msg)
			} else {
				p.HandleAck(ack)
			}
		}
	}
}

func (p *MyPeer) HandleAck(ack *microblock.Ack) {
	logger.Info().Msgf("Received ack message with MB ID %x", ack.MicroblockID)
	//if config.Configuration.MemType == "time" {
	//	r.estimator.AddAck(ack)
	if p.Mempool.IsStable(ack.MicroblockID) {
		return
	}
	if ack.Receiver != util.NodeID(p.Id) {
		// @TODO
		// voteIsVerified, err := crypto.PubVerify(ack.Signature, crypto.IDToByte(ack.MicroblockID), ack.Receiver)
		voteIsVerified := true
		var err error

		if err != nil {
			// log.Warningf("[%v] Error in verifying the signature in ack id: %x", r.ID(), ack.MicroblockID)
			return
		}
		if !voteIsVerified {
			// log.Warningf("[%v] received an ack with invalid signature. vote id: %x", r.ID(), ack.MicroblockID)
			return
		}
	}
	p.Mempool.AddAck(ack)
	found, _ := p.Mempool.FindMicroblock(ack.MicroblockID)
	// @TODO 完成丢失请求部分
	if !found && p.Mempool.IsStable(ack.MicroblockID) {
		missingRequest := microblock.MissingMBRequest{
			RequesterID:   util.NodeID(p.Id),
			MissingMBList: []util.Identifier{ack.MicroblockID},
		}
		msg := &pb.ProtocolMessage{
			SenderId: int32(p.Id),
			Msg: &pb.ProtocolMessage_MissingMicroblockRequest{
				MissingMicroblockRequest: microblock.ToProtoMissingMBRequest(&missingRequest),
			},
		}
		p.SendProtocolMsg(int(ack.Receiver), msg)
	}
}

func (p *MyPeer) HandleMissingMicroblockRequest(mbr *microblock.MissingMBRequest) {
	logger.Info().Msgf("Received mbr message from sender %d with SN %v", mbr.RequesterID, mbr.MissingMBList)
	p.MissingCounts[mbr.RequesterID] += len(mbr.MissingMBList)
	for _, mbid := range mbr.MissingMBList {
		found, mb := p.Mempool.FindMicroblock(mbid)
		if found {
			mb.IsRequested = true
			msg := &pb.ProtocolMessage{
				SenderId: int32(p.Id),
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: microblock.ToProtoMicroBlock(mb),
				},
			}
			p.SendProtocolMsg(int(mbr.RequesterID), msg)
		} else {
			// log.Errorf("[%v] a requested microblock is not found in mempool, id: %x", r.ID(), mbid)
		}
	}
}

// --- 发送端负载均衡部分 --- //
func (p *MyPeer) SenderLoadBalance(protocolMsg *pb.ProtocolMessage) error {
	mb := protocolMsg.Msg.(*pb.ProtocolMessage_Microblock).Microblock
	mb.IsForward = true
	pMsg := &pb.ProtocolMessage{
		SenderId: int32(p.Id),
		Msg: &pb.ProtocolMessage_Microblock{
			Microblock: mb,
		},
	}
	// @TODO
	pick := pickRandomPeer(len(p.PeerConnections), 1, 0)[0]
	logger.Debug().Msgf("[%v] is going to forward a mb to %v", p.Id, pick)
	p.SendProtocolMsg(pick, pMsg)
	return nil
}

func pickRandomPeer(n, d, index int) []int {
	pick := util.RandomPick(n-index, d)
	pickedNode := make([]int, d)
	for i, item := range pick {
		pickedNode[i] = item + 1 + index
	}
	return pickedNode
}

// --- 接收端负载均衡 --- // 
func (p *MyPeer) ReceiverLoadBalance(protocolMsg *pb.ProtocolMessage) error {
	
	return nil
}