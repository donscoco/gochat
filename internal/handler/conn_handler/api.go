package conn_handler

import (
	"encoding/json"
	"errors"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/ws_sender"
	"strconv"
)

type ConnEngine struct{}

const (
	ConnEngineApiSendMsg = "ConnEngine.SendMsg" // todo 所有要发送websockt消息的服务都调用这个rpc接口，后续考虑换成 mq 是不是更好？
)

/*
content:"test"
id:1043
recvId:67
sendId:68
sendTime:1675813467156
type:0


content:"test"
groupId:26
id:225
sendId:68
sendTime:1676111337000
type:0
*/

type SendMsgReq struct {
	Content  string `json:"content"`
	Seq      int64  `json:"seq"` // 消息序列
	RecvId   int    `json:"recvId"`
	GroupId  int    `json:"groupId"`
	SendId   int    `json:"sendId"`
	SendTime int64  `json:"sendTime"`
	Type     int    `json:"type"` // 消息类型

	CmdCode int // 设置cmd code 类型 2 通知下线，3 私聊消息，4 群聊消息
}
type SendMsgReply struct {
	Success bool
	Err     error
}

// 好友列表
func (ce *ConnEngine) SendMsg(req SendMsgReq, reply *SendMsgReply) (err error) {
	//log.Println("rpc recv ", req.Content)

	bl.Info("[SendMsg] recv: ", req.Content)

	ch, ok := ws_sender.GetChan(strconv.Itoa(req.RecvId))
	if !ok {
		bl.Error("")
		reply.Err = errors.New("get websocket fail") // todo
		return
	}

	// 转换到 model 的逻辑，虽然这里全都一样
	wsData := model.SendMsg{
		Content:  req.Content,
		Seq:      req.Seq,
		RecvId:   req.RecvId,
		GroupId:  req.GroupId,
		SendId:   req.SendId,
		SendTime: req.SendTime,
		Type:     req.Type,
	}
	if req.CmdCode == 0 {
		req.CmdCode = 3
	}
	wsMsg := model.WebSocketResponse{
		Cmd:  req.CmdCode, // todo 4 群聊消息
		Data: wsData,
	}
	data, err := json.Marshal(wsMsg)
	if err != nil {
		bl.Errorf("[SendMsg] %+v", err)
	}

	// 序列化一下
	ch <- data
	reply.Success = true
	return nil
}
