package main

import (
	"fmt"
	"log"
	"time"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/mypart/clientstub"
	"github.com/hyperledger-labs/mirbft/mypart/myPeer"
)

func main() {
	// 初始化客户端测试桩
	client := clientStub.NewClientStub(1) // 这里假设客户端 ID 为 1

	// 生成三个客户端请求
	var reqs []*pb.ClientRequest
	for i := 0; i < 3; i++ {
		req := client.GenerateRequest()
		reqs = append(reqs, req)
		log.Printf("Generated request %d: %+v\n", i+1, req)
	}

	// 初始化四个对等节点并彼此相互连接
	var peers []*myPeer.MyPeer
	for i := 1; i <= 4; i++ {
		p, err := myPeer.NewMyPeer(i)
		if err != nil {
			log.Printf("Error creating peer %d: %v\n", i, err)
			return
		}
		peers = append(peers, p)
	}

	// 设置第一个对等节点监听客户端请求并将其广播到其他对等节点
	go func() {
		err := peers[0].ConnectToAllPeers(2, 4)
		if err != nil {
			fmt.Printf("Error connecting to peers: %v\n", err)
			return
		}
		
		// 发送三个请求到其他节点
		for _, req := range reqs {
			err := peers[0].BroadcastRequest(req)
			if err != nil {
				fmt.Printf("Error broadcasting request: %v\n", err)
				return
			}
			time.Sleep(1 * time.Second) // 休眠一秒钟，模拟请求之间的间隔
		}
	}()

	// 阻塞主进程
	select {}
}
