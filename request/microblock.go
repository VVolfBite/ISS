package request

import (
	"encoding/binary"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger-labs/mirbft/crypto"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
	"github.com/kelindar/bitmap"
)

type PendingMicroblock struct {
	Microblock *MicroBlock
	AckMap     map[int32]struct{} // who has sent acks
}
type MicroBlock struct {
	Sn      int
	BucketID        int
	Hash            util.Identifier
	Txns            []*Request
	Timestamp       time.Time
	FutureTimestamp time.Time
	Sender          int32
	IsRequested     bool
	IsForward       bool
	Bitmap          bitmap.Bitmap
	Hops            int
}
type Ack struct {
	Receiver     int32
	BucketID     int
	MicroblockID util.Identifier
	Signature    []byte
}
type PendingBlock struct {
	Payload    *Payload // microblocks that already exist
	MBHashList [][]byte
	MissingMap map[util.Identifier]struct{} // missing list
}
type Payload struct {
	MicroblockList []*MicroBlock
	SigMap         map[util.Identifier]map[int32][]byte
}
type MissingMBRequest struct {
	RequesterID   int32
	BucketID      int32
	Sn    int
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
	hashList = append(hashList, []byte(mb.Timestamp.String()))

	// 计算哈希值并转换成 Identifier 类型
	return util.BytesToIdentifier(crypto.MerkleHashDigests(hashList))
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
	mb.Timestamp = time.Now()
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
		Sn:    int(protoReq.Sn),
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
		Sn:    int32(req.Sn),
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
		Sn:      		int32(mb.Sn),
		Hash:            &pb.Identifier{Value: mb.Hash[:]},
		Txns:            txns,
		Timestamp:       TimeToProtoTimestamp(mb.Timestamp),
		FutureTimestamp: TimeToProtoTimestamp(mb.FutureTimestamp),
		Sender:          int32(mb.Sender),
		IsRequested:     mb.IsRequested,
		IsForward:       mb.IsForward,
		Bitmap:          mb.Bitmap.ToBytes(), // 假设 Bitmap 是一个包含 Bytes() 方法的类型
		Hops:            int32(mb.Hops),
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
		Sn:      		int(protoMb.Sn),
		Hash:            util.BytesToIdentifier(protoMb.Hash.Value),
		Txns:            txns,
		Timestamp:       ProtoTimestampToTime(protoMb.Timestamp),
		FutureTimestamp: ProtoTimestampToTime(protoMb.FutureTimestamp),
		Sender:          int32(protoMb.Sender),
		IsRequested:     protoMb.IsRequested,
		IsForward:       protoMb.IsForward,
		Bitmap:          bitmap.FromBytes(protoMb.Bitmap),
		Hops:            int(protoMb.Hops),
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
