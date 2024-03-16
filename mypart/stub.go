package main

import (
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	clientStub "github.com/hyperledger-labs/mirbft/mypart/clientstub"
	"github.com/hyperledger-labs/mirbft/mypart/myPeer"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	client *clientStub.ClientStub
	peers  []*myPeer.MyPeer
	reqs   []*pb.ClientRequest
)

func Init(reqCount int, peerCount int) {
	// 初始化客户端测试桩

	client = clientStub.NewClientStub(1) // 假设客户端 ID 为 1

	// 生成三个客户端请求
	for i := 0; i < reqCount; i++ {
		req := client.GenerateRequest()
		reqs = append(reqs, req)
		log.Debug().Msgf("Generated request %d: %+v\n", i+1, req)
	}

	// 初始化四个对等节点并彼此相互连接
	for i := 0; i < peerCount; i++ {
		p, err := myPeer.NewMyPeer(i)
		if err != nil {
			log.Printf("Error creating peer %d: %v\n", i, err)
			return
		}
		peers = append(peers, p)
	}
	// 连接对等节点
	for i := 0; i < len(peers); i++ {
		err := peers[i].ConnectToAllPeers(0, peerCount-1)
		if err != nil {
			log.Printf("Error connecting to peers from peer %d: %v\n", i, err)
			return
		}
		// fmt.Printf("%d",len(peers[i].PeerConnections))
	}
}

func testRequestMsg() {
	// 选择第一个对等节点来广播请求
	peer := peers[0]

	// 广播请求
	for _, req := range reqs {
		err := peer.BroadcastRequest(req)
		if err != nil {
			log.Printf("Error broadcasting request: %v\n", err)
			return
		}
		// 休眠一段时间，模拟请求之间的间隔
		time.Sleep(1 * time.Second)
	}
}

func testProtocolMsg() {
	// 创建一个假设的 MicroBlock
	peer := peers[0]
	mb := &pb.MicroBlock{
		ProposalId:      &pb.Identifier{Value: []byte("proposal_id")},
		Hash:            &pb.Identifier{Value: []byte("hash")},
		Sender:          1,
		IsRequested:     true,
		IsForward:       false,
		Timestamp:       &timestamp.Timestamp{Seconds: time.Now().Unix(), Nanos: int32(time.Now().Nanosecond())},
		FutureTimestamp: &timestamp.Timestamp{Seconds: time.Now().Unix(), Nanos: int32(time.Now().Nanosecond())},
		Hops:            1,
	}
	pMsg := &pb.ProtocolMessage{
		SenderId: int32(peer.Id),
		Msg: &pb.ProtocolMessage_Microblock{
			Microblock: mb,
		},
	}
	// 广播 MicroBlock

	// 广播请求
	err := peer.BroadcastProtocolMsg(pMsg)
	if err != nil {
		log.Printf("Error broadcasting request: %v\n", err)
		return
	}
}

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	Init(3, 10)
	testRequestMsg()
	// testProtocolMsg()

	// 阻塞主进程
	select {}
}
