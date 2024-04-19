package request

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger-labs/mirbft/crypto"
	"github.com/hyperledger-labs/mirbft/membership"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
	// logger "github.com/rs/zerolog/log"
)

// 这里引入了新的数据结构


// MicroBlock 用于数据分发的基本分发单元
// @TODO 部分字段尚未启用，如TImeStamp以及Hops
type MicroBlock struct {
	Sn              int
	BucketID        int
	Hash            util.Identifier
	Txns            []*Request
	Sender          int32
	IsRequested     bool
	IsForward       bool
}

// PendingMicroblock 用于存储尚未收集到足够ACK的MB ，ACK存储在ACKBuffer中，所以这里只需要一个快照的Map
type PendingMicroblock struct {
	Microblock *MicroBlock
	AckMap     map[int32]struct{}
}

// Ack用于回复收到了MB，这里比较尴尬的是BucketID字段，因为MB是根据Bucket产生的，不过这并不是一个影响性能和安全的地方
// 不如说ACK的存在使得程序性能提高了，因为使用MBHash——ACK的对应方式，节点只需要验证一次Signature即可，即收到MB时验证一次签名，之后足够的ACK本身保证了MB一定存在以及Req一定合法
type Ack struct {
	Receiver     int32
	BucketID     int
	MicroblockID util.Identifier
	Signature    []byte
}
// PendingBlock 用于存储那些还在等待装载Proposal的Batch，因为Batch是根据一个Bucket提出的，所以PendingBlock对应一个专门的Proposal，
// 目前使用的对应关系是Sn对应 
type PendingBlock struct {
	Payload    *Payload // microblocks that already exist
	MBHashList [][]byte
	MissingMap map[util.Identifier]struct{} // missing list
}
// Payload实际本质是ReqList，他是用于使用Hash进行分发前的本体备份
type Payload struct {
	MicroblockList []*MicroBlock
	SigMap         map[util.Identifier]map[int32][]byte
}
// MissingMBRequest用于召回丢失的MB，这里我们必须处理一个情况，是否会存在即使有足够的ACK证明，但是仍然丢失了MB的情况，导致程序直接卡死？
// FillBatch是另一个影响性能的地方，尤其在真的发生丢包情况时
// 我们需要设计更好的召回方案... 
// 比如再收到不存在的ACK时就可以开始召回了
// 另外应该设计黑名单机制，避免重复请求占用太多带宽
type MissingMBRequest struct {
	RequesterID   int32
	BucketID      int32
	Sn            int
	MissingMBList []util.Identifier
}
type Block struct {
	Payload *Payload
}

func (pd *PendingBlock) CompleteBlock() *Block {
	if len(pd.MissingMap) == 0 {
		return BuildBlock(pd.MBHashList, pd.Payload)
	}
	return nil
}

func (mb *MicroBlock) hash() util.Identifier {
	hashList := make([][]byte, 0)
	for _, tx := range mb.Txns {
		hashList = append(hashList, tx.Digest)
	}
	// 计算哈希值并转换成 Identifier 类型
	return util.BytesToIdentifier(crypto.MerkleHashDigests(hashList))
}

func (ack *Ack) AckVerify() error {
	pk, err := crypto.PublicKeyFromBytes(membership.NodeIdentity(ack.Receiver).PubKey)
	if err != nil {
		return fmt.Errorf("could not verify ack signature: %s", err)
	}
	hash := crypto.Hash(ack.AckDigest())
	err = crypto.CheckSig(hash, pk, ack.Signature)
	if err != nil {
		return fmt.Errorf("could not verify ack signature: %s", err)
	}
	return nil
}

func (ack *Ack) AckDigest() []byte {
	buffer := make([]byte, 0, 4+4+len(ack.MicroblockID))
	receiverBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(receiverBytes, uint32(ack.Receiver))
	bucketIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bucketIDBytes, uint32(ack.BucketID))
	buffer = append(buffer, receiverBytes...)
	buffer = append(buffer, bucketIDBytes...)
	buffer = append(buffer, util.IdentifierToBytes(ack.MicroblockID)...)
	return crypto.Hash(buffer)
	// membership.OwnPrivKey
}

func NewMicroblock(Sn int, txnList []*Request) *MicroBlock {
	mb := new(MicroBlock)
	mb.Sn = Sn
	mb.Txns = txnList
	mb.Hash = mb.hash()
	return mb
}

func NewPendingBlock(MBHashList [][]byte, missingMap map[util.Identifier]struct{}, microBlocks []*MicroBlock) *PendingBlock {
	return &PendingBlock{
		MBHashList: MBHashList,
		MissingMap: missingMap,
		Payload:    &Payload{MicroblockList: microBlocks},
	}
}

func NewPayload(microblockList []*MicroBlock, sigs map[util.Identifier]map[int32][]byte) *Payload {
	return &Payload{
		MicroblockList: microblockList,
		SigMap:         sigs,
	}
}
func (pl *Payload) GenerateHashList() [][]byte {
	hashList := make([][]byte, 0)
	for _, mb := range pl.MicroblockList {
		if mb == nil {
			continue
		}
		hashList = append(hashList, mb.Hash[:])
	}
	return hashList
}

func (pd *PendingBlock) AddMicroblock(mb *MicroBlock) *Block {
	_, exists := pd.MissingMap[mb.Hash]
	if exists {
		pd.Payload.addMicroblock(mb)
		delete(pd.MissingMap, mb.Hash)
	}
	if len(pd.MissingMap) == 0 {
		return BuildBlock(pd.MBHashList, pd.Payload)
	}
	return nil
}

func BuildBlock(MBHashList [][]byte, payload *Payload) *Block {
	return &Block{
		Payload: payload,
	}
}

func (pl *Payload) addMicroblock(mb *MicroBlock) {
	pl.MicroblockList = append(pl.MicroblockList, mb)
}

func FromProtoMissingMBRequest(protoReq *pb.MissingMBRequest) *MissingMBRequest {
	missingMBList := make([]util.Identifier, len(protoReq.MissingMbList))
	for i, id := range protoReq.MissingMbList {
		missingMBList[i] = util.Identifier(FromProtoIdentifier(id))
	}

	return &MissingMBRequest{
		RequesterID:   protoReq.RequesterId,
		Sn:            int(protoReq.Sn),
		MissingMBList: missingMBList,
		BucketID:      protoReq.BucketId,
	}
}

func ToProtoMissingMBRequest(req *MissingMBRequest) *pb.MissingMBRequest {
	protoMissingMBList := make([]*pb.Identifier, len(req.MissingMBList))
	for i, id := range req.MissingMBList {
		protoMissingMBList[i] = ToProtoIdentifier(id)
	}

	return &pb.MissingMBRequest{
		BucketId:      req.BucketID,
		RequesterId:   req.RequesterID,
		Sn:            int32(req.Sn),
		MissingMbList: protoMissingMBList,
	}
}

func FromProtoNodeID(protoID *pb.NodeID) int32 {
	return int32(protoID.GetValue())
}

// Go结构体中的NodeID类型转换为Protobuf消息
func ToProtoNodeID(id int32) *pb.NodeID {
	return &pb.NodeID{
		Value: int32(id),
	}
}

func FromProtoIdentifier(protoID *pb.Identifier) util.Identifier {
	var idBytes [32]byte
	copy(idBytes[:], protoID.GetValue())
	return util.Identifier(idBytes)
}

// Go结构体中的Identifier类型转换为Protobuf消息
func ToProtoIdentifier(id util.Identifier) *pb.Identifier {
	return &pb.Identifier{
		Value: id[:],
	}
}

// 将 Go 中的 MicroBlock 结构体转换为 Protocol Buffers 中的 MicroBlock 消息
func ToProtoMicroBlock(mb *MicroBlock) *pb.MicroBlock {
	var txns []*pb.Request
	for _, txn := range mb.Txns {
		protoTxn := ToProtoRequest(txn) // 将每个请求转换为 Protocol Buffers 中的 Request 消息
		txns = append(txns, protoTxn)
	}

	return &pb.MicroBlock{
		Sn:              int32(mb.Sn),
		Hash:            &pb.Identifier{Value: mb.Hash[:]},
		Txns:            txns,
		Sender:          int32(mb.Sender),
		IsRequested:     mb.IsRequested,
		IsForward:       mb.IsForward,
		BucketId:        int32(mb.BucketID),
	}
}

// 将 Protocol Buffers 中的 MicroBlock 消息转换为 Go 中的 MicroBlock 结构体
func FromProtoMicroBlock(protoMb *pb.MicroBlock) *MicroBlock {
	var txns []*Request
	for _, protoTxn := range protoMb.Txns {
		txn := FromProtoRequest(protoTxn) // 将每个 Protocol Buffers 中的 Request 消息转换为 Go 中的请求
		txns = append(txns, txn)
	}

	return &MicroBlock{
		Sn:              int(protoMb.Sn),
		Hash:            util.BytesToIdentifier(protoMb.Hash.Value),
		Txns:            txns,
		Sender:          int32(protoMb.Sender),
		IsRequested:     protoMb.IsRequested,
		IsForward:       protoMb.IsForward,
		BucketID:        int(protoMb.BucketId),
	}
}

// 将 Go 中的 Ack 结构体转换为 Protocol Buffers 中的 Ack 消息
func ToProtoAck(ack *Ack) *pb.Ack {
	return &pb.Ack{
		Receiver:     int32(ack.Receiver),
		MicroblockId: &pb.Identifier{Value: ack.MicroblockID[:]},
		Signature:    &pb.Signature{Value: ack.Signature},
		BucketId:     int32(ack.BucketID),
	}
}

// 将 Protocol Buffers 中的 Ack 消息转换为 Go 中的 Ack 结构体
func FromProtoAck(protoAck *pb.Ack) *Ack {
	return &Ack{
		Receiver:     int32(protoAck.Receiver),
		MicroblockID: util.BytesToIdentifier(protoAck.MicroblockId.Value),
		Signature:    protoAck.Signature.Value,
		BucketID:     int(protoAck.BucketId),
	}
}

// 将 Go 中的时间类型转换为 Protocol Buffers 中的 Timestamp 类型
func TimeToProtoTimestamp(t time.Time) *timestamp.Timestamp {
	ts, _ := ptypes.TimestampProto(t)
	return ts
}

// 将 Protocol Buffers 中的 Timestamp 类型转换为 Go 中的时间类型
func ProtoTimestampToTime(ts *timestamp.Timestamp) time.Time {
	t, _ := ptypes.Timestamp(ts)
	return t
}

// 将 Go 中的 Request 结构体转换为 Protocol Buffers 中的 Request 消息
func ToProtoRequest(req *Request) *pb.Request {
	return &pb.Request{
		Msg:      req.Msg,
		Digest:   req.Digest,
		Verified: req.Verified,
		InFlight: req.InFlight,
		// 以下两个变量是描述其在本节点的链接情况 因此传输时不再具有意义
	}
}

// 将 Protocol Buffers 中的 Request 消息转换为 Go 中的 Request 结构体
func FromProtoRequest(protoReq *pb.Request) *Request {
	return &Request{
		Msg:      protoReq.Msg,
		Digest:   protoReq.Digest,
		Verified: protoReq.Verified,
		InFlight: protoReq.InFlight,
		// 以下两个变量是描述其在本节点的链接情况 因此传输时不再具有意义
		Next: nil,
		Prev: nil,
	}
}
