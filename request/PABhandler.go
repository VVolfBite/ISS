package request

import (
	// "math/rand"
	// "math/rand"
	"math/rand"
	"time"

	"github.com/hyperledger-labs/mirbft/crypto"
	"github.com/hyperledger-labs/mirbft/membership"
	"github.com/hyperledger-labs/mirbft/messenger"
	pb "github.com/hyperledger-labs/mirbft/protobufs"
	"github.com/hyperledger-labs/mirbft/util"
	logger "github.com/rs/zerolog/log"
)

func HandlePABMsg(protocolMsg *pb.ProtocolMessage) {
	// logger.Info().Msg("PAB IN")
	switch msg := protocolMsg.Msg.(type) {
	case *pb.ProtocolMessage_Microblock:
		// 接收到 Microblock 消息
		// logger.Info().Msg("MB IN")
		mb := msg.Microblock
		rand.Seed(time.Now().UnixNano())
		randomNumber := rand.Float64()
		if randomNumber < 0.01 {
			logger.Info().Msgf("Dropping")
			return
		}
		HandleMicroblock(FromProtoMicroBlock(mb))
		// ack相关内容应当在handler中处理
	case *pb.ProtocolMessage_MicroblockAck:
		// 接收到 Ack 消息
		// logger.Info().Msg("ACK IN")
		ack := msg.MicroblockAck
		HandleAck(FromProtoAck(ack))
	case *pb.ProtocolMessage_MissingMicroblockRequest:
		// 接收到
		// logger.Info().Msg("MRB IN")
		missmbReq := msg.MissingMicroblockRequest
		HandleMissingMicroblockRequest(FromProtoMissingMBRequest(missmbReq))
	default:
		// 其他类型的消息，暂不处理
		logger.Printf("Received unknown message type: %T\n", msg)
	}
}

// bucket部分还不需要实现

func HandleMicroblock(mb *MicroBlock) {
	// if mb.IsRequested{
	// 	logger.Info().Msgf("Rece a mb %x with sn:%d",mb.Hash,mb.Sn)
	// }
	mu.Lock()
	defer mu.Unlock()
	 	
	_, exist := ReceivedMBs[mb.Hash]

	if exist {
		return
	}
	ReceivedMBs[mb.Hash] = struct{}{}
	mb.FutureTimestamp = time.Now()
	sn, exists := MissingMBs[mb.Hash]
	if exists {
		pd, exists := PendingBlockMap[sn]
		logger.Info().Msgf("And its a missing one , a pd with %d to go",len(pd.MissingMap))
		if exists {
			block := pd.AddMicroblock(mb)
			logger.Info().Msgf("Adding to pending block and still got %d to go",len(pd.MissingMap))
			if block != nil {
				// 不能删除这个pd，因为另外一个线程在检索这个东西
				// delete(PendingBlockMap, sn)
				logger.Info().Msgf("A block complete")
				delete(MissingMBs, mb.Hash)
			}
		}
	} else {
		// 不然对方却的也是自己缺少的，加入mempool
		err := Buckets[mb.BucketID].Mempool.AddMicroblock(mb)
		if err != nil {
			logger.Error().Msg("Adding incoming microblock failed")
		}
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
		if mb.IsForward {
			pMsg := &pb.ProtocolMessage{
				SenderId: membership.OwnID,
				Msg: &pb.ProtocolMessage_Microblock{
					Microblock: ToProtoMicroBlock(mb),
				},
			}
			for _, nodeID := range membership.AllNodeIDs() {
				messenger.EnqueueMsg(pMsg, nodeID)
			}
		}
	}
}

func HandleAck(ack *Ack) {
	//if config.Configuration.MemType == "time" {
	//	r.estimator.AddAck(ack)
	// logger.Info().Msgf("Receing ack on %x from %d",ack.MicroblockID,ack.Receiver)
	if Buckets[ack.BucketID].Mempool.IsStable(ack.MicroblockID) {
		return
	}
	if ack.Receiver != membership.OwnID {
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
	Buckets[ack.BucketID].Mempool.AddAck(ack)
	found, _ := Buckets[ack.BucketID].Mempool.FindMicroblock(ack.MicroblockID)
	// @TODO 完成丢失请求部分
	if !found && Buckets[ack.BucketID].Mempool.IsStable(ack.MicroblockID) {
		missingRequest := MissingMBRequest{
			RequesterID:   membership.OwnID,
			MissingMBList: []util.Identifier{ack.MicroblockID},
			Sn:				-1,
			BucketID: 		int32(ack.BucketID),
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

func HandleMissingMicroblockRequest(mbr *MissingMBRequest) {
	// logger.Info().Msgf("Rece a mbr from %d with sn:%d",mbr.RequesterID,mbr.Sn)
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
			logger.Info().Msgf("Fullfil a mb %x from mbr for %d with sn:%d by node: %d",mb.Hash, mbr.RequesterID,mbr.Sn,membership.OwnID)
			messenger.EnqueuePriorityMsg(msg, mbr.RequesterID)
		} else {
			// log.Errorf("[%v] a requested microblock is not found in mempool, id: %x", r.ID(), mbid)
		}
	}
}

// --- 发送端负载均衡部分 --- //
func SenderLoadBalance(protocolMsg *pb.ProtocolMessage) error {
	mb := protocolMsg.Msg.(*pb.ProtocolMessage_Microblock).Microblock
	mb.IsForward = true
	pMsg := &pb.ProtocolMessage{
		SenderId: membership.OwnID,
		Msg: &pb.ProtocolMessage_Microblock{
			Microblock: mb,
		},
	}
	// @TODO
	pick := pickRandomPeer(len(membership.AllNodeIDs()), 1, 0)[0]
	logger.Debug().Msgf("[%v] is going to forward a mb to %v", membership.OwnID, pick)
	messenger.EnqueueMsg(pMsg, int32(pick))
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
func ReceiverLoadBalance(protocolMsg *pb.ProtocolMessage) error {
	
	return nil
}
