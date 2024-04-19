package request

import (
	"math/rand"
	"time"

	"github.com/hyperledger-labs/mirbft/crypto"
	"github.com/hyperledger-labs/mirbft/membership"
	"github.com/hyperledger-labs/mirbft/messenger"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
	logger "github.com/rs/zerolog/log"
)

//这一部分都可以和其他handler合并

func HandlePABMsg(protocolMsg *pb.ProtocolMessage) {
	switch msg := protocolMsg.Msg.(type) {
	case *pb.ProtocolMessage_Microblock:
		mb := msg.Microblock
		HandleMicroblock(FromProtoMicroBlock(mb))
		// ack相关内容应当在handler中处理
	case *pb.ProtocolMessage_MicroblockAck:
		// 接收到 Ack 消息
		ack := msg.MicroblockAck
		HandleAck(FromProtoAck(ack))
	case *pb.ProtocolMessage_MissingMicroblockRequest:
		// 接收到
		missmbReq := msg.MissingMicroblockRequest
		HandleMissingMicroblockRequest(FromProtoMissingMBRequest(missmbReq))
	default:
		// 其他类型的消息，暂不处理
		logger.Printf("Received unknown message type: %T\n", msg)
	}
}

// bucket部分还不需要实现

func HandleMicroblock(mb *MicroBlock) {
	mu.Lock()
	defer mu.Unlock()
	// 检查是否已经收到过
	_, exist := ReceivedMBs[mb.Hash]
	if exist {
		logger.Info().Msgf("Peer %d receive a duplicate mb [%x] , ignoring ... ", membership.OwnID, mb.Hash)
		return
	}

	// 检查是否时MBR请求，处理丢失的MBR
	ReceivedMBs[mb.Hash] = struct{}{}
	sn, exists := MissingMBs[mb.Hash]
	if exists {
		pd, exists := PendingBlockMap[sn]
		logger.Info().Msgf("Peer %d receive a missing mb , a pending block with %d to go", membership.OwnID, len(pd.MissingMap))
		if exists {
			block := pd.AddMicroblock(mb)
			logger.Info().Msgf("Adding to pending block and still got %d to go", len(pd.MissingMap))
			if block != nil {
				// 不能删除这个pd，因为另外一个线程在检索这个东西
				// delete(PendingBlockMap, sn)
				logger.Info().Msgf("A pending block complete to a block")
				delete(MissingMBs, mb.Hash)
			}
		}
	} else {
		// 否则就是正常的对方生成了MB
		err := Buckets[mb.BucketID].Mempool.AddMicroblock(mb)
		if err != nil {
			logger.Error().Msg("Adding incoming microblock failed")
		}
		// 回复ACK
		if !mb.IsRequested {
			ack := &Ack{
				Receiver:     membership.OwnID,
				MicroblockID: mb.Hash,
				BucketID:     mb.BucketID,
			}
			sk, err := crypto.PrivateKeyFromBytes(membership.OwnPrivKey)
			if err != nil {
				logger.Error().Any("Err", err).Msg("Signing ack failed")
			}
			hash := crypto.Hash(ack.AckDigest())

			signature, err := crypto.Sign(hash, sk)
			if err != nil {
				logger.Error().Any("Err", err).Msg("Signing ack failed")
			}
			ack.Signature = signature
			msg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_MicroblockAck{
					MicroblockAck: ToProtoAck(ack),
				},
			}
			if mb.Sender != membership.OwnID {
				messenger.EnqueueMsg(msg, mb.Sender)
			} else {
				HandleAck(ack)
			}
		}
		// 帮助对方转发
		if mb.IsForward {
			logger.Info().Msgf("Peer %v is going to forward a mb %x for %d ", membership.OwnID, mb.Hash, mb.Sender)
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(mb),
				},
			}
			for _, nodeID := range membership.AllNodeIDs() {
				if nodeID == membership.OwnID{
					continue
				}
				messenger.EnqueueMsg(pMsg, nodeID)
			}
		}
	}
}

func HandleAck(ack *Ack) {
	if Buckets[ack.BucketID].Mempool.IsStable(ack.MicroblockID) {
		return
	}
	if ack.AckVerify() != nil {
		logger.Info().Msgf("Warning: wrong ack received!")
		return
	}
	Buckets[ack.BucketID].Mempool.AddAck(ack)
	found, _ := Buckets[ack.BucketID].Mempool.FindMicroblock(ack.MicroblockID)
	// @TODO 完成丢失请求部分
	if !found && Buckets[ack.BucketID].Mempool.IsStable(ack.MicroblockID) {
		missingRequest := MissingMBRequest{
			RequesterID:   membership.OwnID,
			MissingMBList: []util.Identifier{ack.MicroblockID},
			Sn:            -1,
			BucketID:      int32(ack.BucketID),
		}
		msg := &pb.ProtocolMessage{
			SenderId: membership.OwnID,
			Msg: &pb.ProtocolMessage_MissingMicroblockRequest{
				MissingMicroblockRequest: ToProtoMissingMBRequest(&missingRequest),
			},
		}
		messenger.EnqueueMsg(msg, ack.Receiver)
	}
}

// 丢弃MB的本质是增加了一轮信息时延，不到万不得已，不应当允许节点肆意丢弃MB，
func HandleMissingMicroblockRequest(mbr *MissingMBRequest) {
	MissingCounts[mbr.RequesterID] += len(mbr.MissingMBList)
	for _, mbid := range mbr.MissingMBList {
		found, mb := Buckets[mbr.BucketID].Mempool.FindMicroblock(mbid)
		if found {
			mb.IsRequested = true
			mb.Sn = mbr.Sn
			msg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(mb),
				},
			}
			logger.Debug().Msgf("Fullfil a mb %x from mbr for %d with sn:%d by node: %d", mb.Hash, mbr.RequesterID, mbr.Sn, membership.OwnID)
			messenger.EnqueuePriorityMsg(msg, mbr.RequesterID)
		} else {
			// log.Errorf("[%v] a requested microblock is not found in mempool, id: %x", r.ID(), mbid)
		}
	}
}

func pickRandomNode() int32 {
	allNodes := membership.AllNodeIDs()
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(allNodes))
	return allNodes[randomIndex]
}
