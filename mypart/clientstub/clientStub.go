package clientStub

import (
	"fmt"
	"github.com/hyperledger-labs/mirbft/protobufs"
	"math/rand"
	"time"
	// _logger "github.com/rs/zerolog/log"
)

// ClientStub 定义一个客户端测试桩
type ClientStub struct {
	ID      int32 // 客户端ID
	counter int32 // 请求计数器
}

// NewClientStub 初始化并返回一个客户端测试桩实例
func NewClientStub(id int32) *ClientStub {
	return &ClientStub{
		ID:      id,
		counter: 0, // 初始化计数器
	}
}

// GenerateRequest 生成一个随机的客户端请求
func (c *ClientStub) GenerateRequest() *protobufs.ClientRequest {
	// 使用当前时间作为随机数种子
	rand.Seed(time.Now().UnixNano())

	// 递增计数器
	c.counter++

	// 构建一个请求
	req := &protobufs.ClientRequest{
		RequestId: &protobufs.RequestID{
			ClientId: c.ID,
			ClientSn: c.counter,
		},
		Payload:   []byte(fmt.Sprintf("payload for client %d, request %d", c.ID, c.counter)),
		Pubkey:    nil, // 可以根据需要设置公钥
		Signature: nil, // 可以根据需要设置签名
	}

	return req
}

func main() {
	// 初始化测试桩
	clientStub := NewClientStub(20) // 这里假设 ID 为 20

	// 生成一个随机的请求
	req := clientStub.GenerateRequest()

	// 输出请求信息
	fmt.Printf("Generated request: %+v\n", req)
}
