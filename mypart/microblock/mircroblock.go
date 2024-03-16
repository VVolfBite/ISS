package microblock

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger-labs/mirbft/crypto"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/request"
	"github.com/hyperledger-labs/mirbft/util"
	"github.com/kelindar/bitmap"
	"github.com/rs/zerolog/log"
)

type PendingMicroblock struct {
	Microblock *MicroBlock
	AckMap     map[util.NodeID]struct{} // who has sent acks
}
type MicroBlock struct {
	ProposalID      util.Identifier
	Hash            util.Identifier
	Txns            []*request.Request
	Timestamp       time.Time
	FutureTimestamp time.Time
	Sender          util.NodeID
	IsRequested     bool
	IsForward       bool
	Bitmap          bitmap.Bitmap
	Hops            int
}

type PendingBlock struct {
	Payload    *Payload // microblocks that already exist
	Proposal   *Proposal
	MissingMap map[util.Identifier]struct{} // missing list
}

type Ack struct {
	Receiver     util.NodeID
	MicroblockID util.Identifier
	util.Signature
}

type Payload struct {
	MicroblockList []*MicroBlock
	SigMap         map[util.Identifier]map[util.NodeID]util.Signature
}

type Proposal struct {
	HashList []util.Identifier
}

type MissingMBRequest struct {
	RequesterID   util.NodeID
	ProposalID    util.Identifier
	MissingMBList []util.Identifier
}

type Block struct {
	Payload *Payload
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

func NewMicroblock(proposalID util.Identifier, txnList []*request.Request) *MicroBlock {
	mb := new(MicroBlock)
	mb.ProposalID = proposalID
	mb.Txns = txnList
	mb.Timestamp = time.Now()
	mb.Hash = mb.hash()
	return mb
}

func NewPendingBlock(proposal *Proposal, missingMap map[util.Identifier]struct{}, microBlocks []*MicroBlock) *PendingBlock {
	return &PendingBlock{
		Proposal:   proposal,
		MissingMap: missingMap,
		Payload:    &Payload{MicroblockList: microBlocks},
	}
}

func NewPayload(microblockList []*MicroBlock, sigs map[util.Identifier]map[util.NodeID]util.Signature) *Payload {
	return &Payload{
		MicroblockList: microblockList,
		SigMap:         sigs,
	}
}

func (pd *PendingBlock) AddMicroblock(mb *MicroBlock) *Block {
	_, exists := pd.MissingMap[mb.Hash]
	if exists {
		pd.Payload.addMicroblock(mb)
		delete(pd.MissingMap, mb.Hash)
	}
	if len(pd.MissingMap) == 0 {
		return BuildBlock(pd.Proposal, pd.Payload)
	}
	return nil
}

func BuildBlock(proposal *Proposal, payload *Payload) *Block {
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
		RequesterID:   util.NodeID(FromProtoNodeID(protoReq.RequesterId)),
		ProposalID:    util.Identifier(FromProtoIdentifier(protoReq.ProposalId)),
		MissingMBList: missingMBList,
	}
}

func ToProtoMissingMBRequest(req *MissingMBRequest) *pb.MissingMBRequest {
	protoMissingMBList := make([]*pb.Identifier, len(req.MissingMBList))
	for i, id := range req.MissingMBList {
		protoMissingMBList[i] = ToProtoIdentifier(id)
	}

	return &pb.MissingMBRequest{
		RequesterId:   ToProtoNodeID(req.RequesterID),
		ProposalId:    ToProtoIdentifier(req.ProposalID),
		MissingMbList: protoMissingMBList,
	}
}

func FromProtoNodeID(protoID *pb.NodeID) util.NodeID {
	return util.NodeID(protoID.GetValue())
}

// Go结构体中的NodeID类型转换为Protobuf消息
func ToProtoNodeID(id util.NodeID) *pb.NodeID {
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
	log.Printf("%+v", mb)
	var txns []*pb.Request
	for _, txn := range mb.Txns {
		protoTxn := ToProtoRequest(txn) // 将每个请求转换为 Protocol Buffers 中的 Request 消息
		txns = append(txns, protoTxn)
	}

	return &pb.MicroBlock{
		ProposalId:      &pb.Identifier{Value: mb.ProposalID[:]},
		Hash:            &pb.Identifier{Value: mb.Hash[:]},
		Txns:            txns,
		Timestamp:       TimeToProtoTimestamp(mb.Timestamp),
		FutureTimestamp: TimeToProtoTimestamp(mb.FutureTimestamp),
		Sender:          int32(mb.Sender),
		IsRequested:     mb.IsRequested,
		IsForward:       mb.IsForward,
		Bitmap:          mb.Bitmap.ToBytes(), // 假设 Bitmap 是一个包含 Bytes() 方法的类型
		Hops:            int32(mb.Hops),
	}
}

// 将 Protocol Buffers 中的 MicroBlock 消息转换为 Go 中的 MicroBlock 结构体
func FromProtoMicroBlock(protoMb *pb.MicroBlock) *MicroBlock {
	var txns []*request.Request
	for _, protoTxn := range protoMb.Txns {
		txn := FromProtoRequest(protoTxn) // 将每个 Protocol Buffers 中的 Request 消息转换为 Go 中的请求
		txns = append(txns, txn)
	}

	return &MicroBlock{
		ProposalID:      util.BytesToIdentifier(protoMb.ProposalId.Value),
		Hash:            util.BytesToIdentifier(protoMb.Hash.Value),
		Txns:            txns,
		Timestamp:       ProtoTimestampToTime(protoMb.Timestamp),
		FutureTimestamp: ProtoTimestampToTime(protoMb.FutureTimestamp),
		Sender:          util.NodeID(protoMb.Sender),
		IsRequested:     protoMb.IsRequested,
		IsForward:       protoMb.IsForward,
		Bitmap:          bitmap.FromBytes(protoMb.Bitmap),
		Hops:            int(protoMb.Hops),
	}
}

// 将 Go 中的 Ack 结构体转换为 Protocol Buffers 中的 Ack 消息
func ToProtoAck(ack *Ack) *pb.Ack {
	return &pb.Ack{
		Receiver:     int32(ack.Receiver),
		MicroblockId: &pb.Identifier{Value: ack.MicroblockID[:]},
		Signature:    &pb.Signature{Value: ack.Signature},
	}
}

// 将 Protocol Buffers 中的 Ack 消息转换为 Go 中的 Ack 结构体
func FromProtoAck(protoAck *pb.Ack) *Ack {
	return &Ack{
		Receiver:     util.NodeID(protoAck.Receiver),
		MicroblockID: util.BytesToIdentifier(protoAck.MicroblockId.Value),
		Signature:    protoAck.Signature.Value,
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
func ToProtoRequest(req *request.Request) *pb.Request {
	return &pb.Request{
		Msg:      req.Msg,
		Digest:   req.Digest,
		Verified: req.Verified,
		InFlight: req.InFlight,
		// 以下两个变量是描述其在本节点的链接情况 因此传输时不再具有意义
	}
}

// 将 Protocol Buffers 中的 Request 消息转换为 Go 中的 Request 结构体
func FromProtoRequest(protoReq *pb.Request) *request.Request {
	return &request.Request{
		Msg:      protoReq.Msg,
		Digest:   protoReq.Digest,
		Verified: protoReq.Verified,
		InFlight: protoReq.InFlight,
		// 以下两个变量是描述其在本节点的链接情况 因此传输时不再具有意义
		Next: nil,
		Prev: nil,
	}
}
