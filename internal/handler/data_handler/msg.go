package data_handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/donscoco/gochat/internal/base/kafka/mkafka"
	"github.com/donscoco/gochat/internal/base/mongodb"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/dao/mgo"
	"github.com/donscoco/gochat/internal/handler/conn_handler"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/conn_manager"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/donscoco/gochat/internal/service/msg_service/msg_producer"
	"github.com/donscoco/gochat/pkg/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
)

const (
	DataEngineApiSendMsg = "DataEngine.SendMsg"
	DataEngineApiPullMsg = "DataEngine.PullMsg" // 获取未读消息
	DataEngineApiHistory = "DataEngine.History" // 查看历史消息
)

type SendMsgReq struct {
	Content  string `json:"content"`
	Seq      int64  `json:"seq"`     // 消息序列
	RecvId   int    `json:"recvId"`  // 在发送私聊消息的时候会有
	GroupId  int    `json:"groupId"` // 在发送群消息的时候会有
	SendId   int    `json:"sendId"`
	SendTime int64  `json:"sendTime"`
	Type     int    `json:"type"` // 消息类型
}
type SendMsgReply struct {
	Success bool
}

func (de *DataEngine) SendMsg(req SendMsgReq, reply *SendMsgReply) (err error) {
	data, _ := json.Marshal(req)
	log.Println(string(data))

	// 消息类型数据直接发到队列去处理

	topic := "domark-test"
	msg := mkafka.NewMessage(topic, data)
	msg_producer.MsgProducer.Input() <- msg

	reply.Success = true

	return
}

type PullMsgReq struct {
	UserId int `json:"userId"` // 在发送私聊消息的时候会有

	//Seq      int64  `json:"seq"`     // 消息序列
	IsGroup bool
}
type PullMsgReply struct {
	Success bool
	Err     error
}

func (de *DataEngine) PullMsg(req PullMsgReq, reply *PullMsgReply) (err error) {

	// 获取 用户的所有 conversation，查看conversation[] 中的read_offset
	// 将比 read offset 还要大的，recvId 是 当前用户的 消息从mongo中拿出来（限制最多100条）
	// 更新最新的read offset
	// 消息时间排序（小到大）
	// 将消息发送到 对应的channel

	tx, err := gorm.GetGormPool("default")
	if err != nil {
		// todo log
		bl.Error("[PullMsg] get gorm fail")
		reply.Err = err
		return
	}

	tx = tx.Begin()
	list := make([]dao.Conversation, 0, 16)
	// 拿到本用户的 read_offset 信息
	// 分私聊群聊
	if req.IsGroup {
		err = tx.Raw("SELECT * FROM chat_conversation "+
			"WHERE user_id = ? AND group_id <> 0 for update", req.UserId).Scan(&list).Error // 拿到所有群聊会话
		if err != nil {
			// todo
			tx.Rollback()
			return
		}

	} else {
		err = tx.Raw("SELECT * FROM chat_conversation "+
			"WHERE user_id = ? AND group_id = 0 for update", req.UserId).Scan(&list).Error // 拿到所有私聊会话
		if err != nil {
			// todo
			tx.Rollback()
			return
		}
	}

	var convMsg = make(map[dao.Conversation][]mgo.MsgRecord) // 放每个会话的消息

	for _, c := range list {
		records := make([]mgo.MsgRecord, 0, 16)
		// 分私聊群聊
		if req.IsGroup {
			records, err = findGroupMsg(&c)
			if err != nil {
				// todo
				tx.Rollback()
				return
			}

		} else {
			records, err = findPrivateMsg(&c)
			if err != nil {
				// todo
				tx.Rollback()
				return
			}
		}

		if len(records) > 0 {
			convMsg[c] = records

			sort.Sort(MsgRecordSlice(records)) // 排序

			// todo 更新 conversation 的 read offset
			err = tx.Model(c).Update("read_offset", records[len(records)-1].Seq).Error // 使用最新的 seq 作为 read_offset
			if err != nil {
				// todo
				tx.Rollback()
				return err

			}
		}
	}

	tx.Commit()
	reply.Success = true

	// todo 后续封装一下模块，第一版没什么内容，先放这里了

	if req.IsGroup {
		for _, records := range convMsg {
			for _, msg := range records {
				// 发送给websocket
				sendReply := conn_handler.SendMsgReply{}
				err = client_service.DefaultConnEngineCli.CallSpecificAddr(
					conn_handler.ConnEngineApiSendMsg,
					conn_handler.SendMsgReq{
						CmdCode: model.CmdGroup,
						Content: msg.Content,
						Seq:     msg.Seq,
						//RecvId:   msg.RecvId,
						RecvId:   req.UserId, // recvId 是要发送给的websockt。这个接口是发送给自己
						GroupId:  msg.GroupId,
						SendId:   msg.SendId, // 群组里面谁发送的这条信息
						SendTime: msg.SendTime,
						Type:     msg.Type,
					},
					&sendReply,
					conn_manager.Where(req.UserId), // 指定 fd 所在的机器提供服务
				)
				if err != nil { // todo
					return err
				}
				if reply.Success != true {
					log.Println("send msg to ws fail") // todo log
					return errors.New("send msg to ws fail")
				}
			}
		}
	} else {
		for _, records := range convMsg {
			for _, msg := range records {
				// 发送给websocket
				sendReply := conn_handler.SendMsgReply{}
				err = client_service.DefaultConnEngineCli.CallSpecificAddr(
					conn_handler.ConnEngineApiSendMsg,
					conn_handler.SendMsgReq{
						CmdCode: model.CmdPrivate,
						Content: msg.Content,
						Seq:     msg.Seq,
						RecvId:  msg.RecvId,
						//GroupId:  msg.GroupId,
						SendId:   msg.SendId,
						SendTime: msg.SendTime,
						Type:     msg.Type,
					},
					&sendReply,
					conn_manager.Where(req.UserId), // 指定 fd 所在的机器提供服务
				)
				if err != nil { // todo
					return err
				}
				if reply.Success != true {
					log.Println("send msg to ws fail") // todo log
					return errors.New("send msg to ws fail")
				}
			}
		}
	}

	return
}

// todo 封装到模块中
func findPrivateMsg(c *dao.Conversation) (records []mgo.MsgRecord, err error) {

	filter := bson.D{
		{"_id",
			bson.D{
				{"$gt", c.ReadOffset},
			},
		},
		{"recvId", c.UserId},
		//{"groupId", c.GroupId},
	}

	cursor, err := mongodb.DefaultMgoDB.Database("gochat").Collection("Msg").
		Find(context.Background(), filter, options.Find().SetLimit(100))
	if err != nil { //todo
		log.Fatalln(err)
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	records = make([]mgo.MsgRecord, 0, 16)
	err = cursor.All(context.TODO(), &records) // 我们要查的数量不多，直接查就行
	if err != nil {
		log.Fatal(err)
	}

	return
}
func findGroupMsg(c *dao.Conversation) (mrs []mgo.MsgRecord, err error) {

	filter := bson.D{
		{"_id",
			bson.D{
				{"$gt", c.ReadOffset},
			},
		},
		{"groupId", c.GroupId},
	}

	cursor, err := mongodb.DefaultMgoDB.Database("gochat").Collection("Msg").
		Find(context.Background(), filter, options.Find().SetLimit(100))
	if err != nil { //todo
		log.Fatalln(err)
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	mrs = make([]mgo.MsgRecord, 0, 16)
	err = cursor.All(context.TODO(), &mrs) // 我们要查的数量不多，直接查就行
	if err != nil {
		log.Fatal(err)
	}

	return
}
func findPrivateHistory(uid, fid, gid, page, size int64) (records []mgo.MsgRecord, err error) {

	filter := bson.D{
		{"$or",
			bson.A{
				bson.M{"recvId": uid, "sendId": fid},
				bson.M{"sendId": uid, "recvId": fid},
			},
		},
	}

	skip := (page - 1) * size

	sort := bson.M{ // 这里不要用时间排序，因为在多个节点的服务中，可能是先生成seq后，另一个及诶单生成seq，时间，然后才到本节点生成时间
		"seq": -1, // 从大到小排序
	}

	cursor, err := mongodb.DefaultMgoDB.Database("gochat").Collection("Msg").
		Find(context.Background(), filter, options.Find().SetSort(sort).SetSkip(skip).SetLimit(size))
	if err != nil { //todo
		log.Fatalln(err)
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	records = make([]mgo.MsgRecord, 0, 16)
	err = cursor.All(context.TODO(), &records) // 我们要查的数量不多，直接查就行
	if err != nil {
		log.Fatal(err)
	}

	return
}
func findGroupHistory(uid, fid, gid, page, size int64) (records []mgo.MsgRecord, err error) {

	filter := bson.M{
		"groupId": gid,
	}
	//filter := bson.D{
	//	{"$or",
	//		bson.A{
	//			bson.M{"recvId": uid, "sendId": fid},
	//			bson.M{"sendId": uid, "recvId": fid},
	//		},
	//	},
	//}

	skip := (page - 1) * size

	sort := bson.M{ // 这里不要用时间排序，因为在多个节点的服务中，可能是先生成seq后，另一个及诶单生成seq，时间，然后才到本节点生成时间
		"seq": -1, // 从大到小排序
	}

	cursor, err := mongodb.DefaultMgoDB.Database("gochat").Collection("Msg").
		Find(context.Background(), filter, options.Find().SetSort(sort).SetSkip(skip).SetLimit(size))
	if err != nil { //todo
		log.Fatalln(err)
	}
	//延迟关闭游标
	defer func() {
		if err = cursor.Close(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	records = make([]mgo.MsgRecord, 0, 16)
	err = cursor.All(context.TODO(), &records) // 我们要查的数量不多，直接查就行
	if err != nil {
		log.Fatal(err)
	}

	return
}

// 按照序列号去排序
type MsgRecordSlice []mgo.MsgRecord

func (x MsgRecordSlice) Len() int           { return len(x) }
func (x MsgRecordSlice) Less(i, j int) bool { return x[i].Seq < x[j].Seq }
func (x MsgRecordSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// /////// history ////////////
type HistoryMsgReq struct {
	Page   int64
	Size   int64
	UserId int64

	FriendId int64 // 查询私聊的id
	GroupId  int64 // 查询群聊的id
}
type HistoryMsgReply struct {
	Records []mgo.MsgRecord
	Success bool
	Err     error
}

func (de *DataEngine) History(req HistoryMsgReq, reply *HistoryMsgReply) (err error) {

	if req.GroupId == 0 {
		reply.Records, err = findPrivateHistory(req.UserId, req.FriendId, req.GroupId, req.Page, req.Size)
		if err != nil {
			reply.Err = err
		}
		return nil
	} else {
		reply.Records, err = findGroupHistory(req.UserId, req.FriendId, req.GroupId, req.Page, req.Size)
		if err != nil {
			reply.Err = err
		}
		return nil
	}

}
