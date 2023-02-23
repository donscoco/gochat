package model

type MsgSendInput struct {
	Content string `json:"content"` // 消息内容
	RecvId  int    `json:"recvId"`  // 消息接收者 // 私聊
	GroupId int    `json:"groupId"` // 消息接收者 // 群聊
	Type    int    `json:"type"`    // 消息类型

	// 上面3个是主要的。一下是一些消息类型发送的数据结构，但是不需要用到。如 sendTime，肯定是按照服务器收到的为准，服务器自己生成时间
	Id         int    `json:"id"`
	FileId     int64  `json:"fileId"`
	SendId     int    `json:"sendId"`
	SendTime   int64  `json:"sendTime"`
	SelfSend   bool   `json:"selfSend"`
	LoadStatus string `json:"loading"`
}

type MsgSendGroupInput struct {
	Content string `json:"content"` // 消息内容
	GroupId int    `json:"groupId"` // 消息接收者
	Type    int    `json:"type"`    // 消息类型
}
