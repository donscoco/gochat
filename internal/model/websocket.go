package model

// 存放发给websocket 的通用结构体

//////   存放websocket 的通用 模版 //////
/*
{
	"cmd": int
	"data":{}
}
*/

const (
	CmdHeartbeat = 1 // 心跳
	CmdLogin     = 2 // 有另一个地方登陆了。发送cmd 2 通知其下线
	CmdPrivate   = 3 // 私聊消息
	CmdGroup     = 4 // 群聊消息
)

type WebSocketResponse struct {
	Cmd  int         `json:"cmd"` // 做成const
	Data interface{} `json:"data"`
}

type SendMsg struct {
	Content  string `json:"content"`
	Seq      int64  `json:"seq"` // 消息序列
	RecvId   int    `json:"recvId"`
	GroupId  int    `json:"groupId"`
	SendId   int    `json:"sendId"`
	SendTime int64  `json:"sendTime"`
	Type     int    `json:"type"` // 消息类型
}
