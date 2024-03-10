package myPeer

import (
	"fmt"
	"log"
	"net"
	"context"
	"github.com/hyperledger-labs/mirbft/mypart/mempool"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"google.golang.org/grpc"
)

// MyPeer 结构体作为 gRPC 服务器的实现类型
type MyPeer struct {
	Id                int
	Mempool           *mempool.MemPool
	PeerConnections   map[int]*grpc.ClientConn // 将值类型改为指针类型
}

func NewMyPeer(Id int) (*MyPeer, error) {
	peer := &MyPeer{
		Id:                Id,
		Mempool:           mempool.NewMemPool(300,300,3000),
		PeerConnections:   make(map[int]*grpc.ClientConn), // 将值类型改为指针类型
	}

	server := grpc.NewServer()
	pb.RegisterMessengerServer(server, peer)
	address := fmt.Sprintf("localhost:%d", 50000+Id)
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return nil, err
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	return peer, nil
}

// ConnectToAllPeers 尝试连接从 start 到 end 的所有对等节点
func (p *MyPeer) ConnectToAllPeers(start int, end int) error {
	for Id := start; Id <= end; Id++ {
		if err := p.ConnectToPeer(Id); err != nil {
			log.Printf("Failed to connect to peer %d: %v", Id, err)
		}
	}
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
	log.Printf("Connecting to peer: %d\n", Id)
	return nil
}

// Listen 实现了 Listen RPC 方法
func (p *MyPeer) Listen(stream pb.Messenger_ListenServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}

		log.Printf("Received message from sender %d with SN %d\n", msg.GetSenderId(), msg.GetSn())
		// 假设你有一个名为 ack 的 BandwIdthTestAck 消息
		ack := &pb.BandwidthTestAck{}
		if err := stream.Send(ack); err != nil {
			return err
		}
	}
}

func (p *MyPeer) Request(stream pb.Messenger_RequestServer) error {
    for {
        req, err := stream.Recv()
        if err != nil {
            return err
        }

        log.Printf("Received client request: %v\n", req)
		p.HandleRequest(req)
        
        // 创建 ClientResponse 消息
        resp := &pb.ClientResponse{
            ClientSn:  req.RequestId.ClientSn,
            // 不设置 order_sn
        }

        // 发送 ClientResponse 响应
        if err := stream.Send(resp); err != nil {
            return err
        }
    }
}

// Buckets 实现了 Buckets RPC 方法
func (p *MyPeer) Buckets(stream pb.Messenger_BucketsServer) error {
	// 实现 Buckets RPC 方法的逻辑
	return nil
}


func (p *MyPeer) SendRequest(Id int, req *pb.ClientRequest) error {
	conn, ok := p.PeerConnections[Id]
	if !ok {
		return fmt.Errorf("peer %d is not connected", Id)
	}

	client := pb.NewMessengerClient(conn)
	stream, err := client.Request(context.Background())
	if err != nil {
		return err
	}

	// 发送请求
	if err := stream.Send(req); err != nil {
		return err
	}

	// 接收响应
	resp, err := stream.Recv()
	if err != nil {
		return err
	}

	log.Printf("Received response from peer %d: %v\n", Id, resp)
	return nil
}

func (p *MyPeer) BroadcastRequest(req *pb.ClientRequest) error {
	for id := range p.PeerConnections {
		err := p.SendRequest(id, req)
		if err != nil {
			log.Printf("Error sending request to peer %d: %v", id, err)
		}
	}
	return nil
}


// HandleRequest 处理请求，简单地调用 Mempool 的 ReceiveTxn 方法
func (p *MyPeer) HandleRequest(req *pb.ClientRequest) {
	// 将请求添加到内存池
	newMb, microBlock := p.Mempool.AddTxn(req)
	if newMb {
		p.Mempool.AddMicroblock(microBlock)
		log.Printf("Txn added to peer %d's mempool with a new mb return.\n", p.Id)
	} else {		
		log.Printf("Txn added to peer %d's mempool\n", p.Id)	
	}
	log.Printf("Peer %d's mempool status:\n", p.Id)
	log.Printf("Total transactions: %d\n", p.Mempool.TotalTx())
	log.Printf("Remaining transactions: %d\n", p.Mempool.RemainingTx())
	log.Printf("Total microblocks: %d\n", p.Mempool.TotalMB())
	log.Printf("Remaining microblocks: %d\n", p.Mempool.RemainingMB())

	// 打印内存池状态

}
