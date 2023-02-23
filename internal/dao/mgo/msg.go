package mgo

type MsgRecord struct {
	Id       int64  `json:"-" bson:"_id"`
	Content  string `json:"content" bson:"content"` //
	Seq      int64  `json:"seq" bson:"seq"`         // 消息序列
	RecvId   int    `json:"recvId" bson:"recvId"`
	GroupId  int    `json:"groupId" bson:"groupId"`
	SendId   int    `json:"sendId" bson:"sendId"`
	SendTime int64  `json:"sendTime" bson:"sendTime"`
	Type     int    `json:"type" bson:"type"` // 消息类型
}

type MsgGroupRecord struct {
	Id       int64  `bson:"_id"`
	Content  string `json:"content" bson:"content"` //
	Seq      int64  `json:"seq" bson:"seq"`         // 消息序列
	RecvId   int    `json:"recvId" bson:"recvId"`
	SendId   int    `json:"sendId" bson:"sendId"`
	GroupId  int    `json:"groupId" bson:"groupId"`
	SendTime int64  `json:"sendTime" bson:"sendTime"`
	Type     int    `json:"type" bson:"type"` // 消息类型
}
