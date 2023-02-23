package msg_process

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/donscoco/gochat/internal/base/kafka/mkafka"
	"github.com/donscoco/gochat/internal/base/mongodb"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/dao/mgo"
	"github.com/donscoco/gochat/internal/handler/conn_handler"
	"github.com/donscoco/gochat/internal/handler/data_handler"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/conn_manager"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/donscoco/gochat/pkg/gorm"
	"log"
	"time"
)

func MsgHandleFunc(kc *mkafka.KafkaConsumer, msg *sarama.ConsumerMessage) (err error) {

	// 分 群消息和 私聊消息
	log.Printf("kafka handler recv : %s ,%+v \n", string(msg.Value), msg)
	// 解析 msg
	// 写入 db
	//     保存进mongodb
	//     开启mysql 事务，更新conversation表，如果没有conversation项就创建一个，记录 write_offset
	// 删除缓存
	//
	// 如果接受者在线，需要更新 接受者的read_offset 并发送websocket （要区分群聊和私聊）
	// 查询 uid 在哪个 addr
	// 调用 conn_rpc

	params := &data_handler.SendMsgReq{}
	err = json.Unmarshal(msg.Value, params)
	if err != nil {
		return err

	} // todo

	/* DE逻辑开始 */
	record := mgo.MsgRecord{
		Id:       params.Seq,     // seq 是全局唯一的
		Content:  params.Content, // 消息内容，todo 是否改解析后展示在 mongodb 中方便查看比较好？
		Seq:      params.Seq,
		RecvId:   params.RecvId,
		GroupId:  params.GroupId,
		SendId:   params.SendId,
		SendTime: params.SendTime,
		Type:     params.Type,
	}
	insertResult, err := mongodb.DefaultMgoDB.Database("gochat").Collection("Msg").
		InsertOne(context.TODO(), record)
	if err != nil {
		// todo
		bl.Error("[msg_process] insert Msg to MongoDB fail ", insertResult)
	}
	if insertResult != nil { // todo rm ,for debug
		bl.Debugf("[msg_process] msg:%d", insertResult.InsertedID.(int64))
	}

	// 开始分别处理

	var isGroupMsg bool
	if params.RecvId != 0 {
		isGroupMsg = false
	}
	if params.GroupId != 0 {
		isGroupMsg = true
	}

	// 被开启会话的一方判断是否有会话并创建的逻辑放到 需要修改 read_offset 的地方写 是否比较好？
	// todo 改成在这里创建被开启会话方，（因为这里是异步队列），分私聊和群组，群聊的话还要查询组的成员信息（新加成员需要在添加会话）
	if isGroupMsg {
		err = Group(params, kc, msg)
	} else {
		err = Private(params, kc, msg)
	}

	if err != nil {
		// todo
		return err
	}

	return
}
func Group(params *data_handler.SendMsgReq, kc *mkafka.KafkaConsumer, msg *sarama.ConsumerMessage) (err error) {

	// 开启事务
	// 获取groupid的所有成员
	// 检查每个成员是否有会话并创建会话
	// 检查每个成员是否在线，对于在线的要更新read_offset并发送websocket

	var onLineMap = make(map[int]bool)
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return err
	}

	tx = tx.Begin()

	// 将发送者的会话创建出来并修改
	search1 := dao.Conversation{} // for write_offset
	tx.Raw("SELECT * FROM chat_conversation "+
		"WHERE user_id = ? AND recv_id = ? AND group_id = ? for update", params.SendId, params.RecvId, params.GroupId).Scan(&search1)
	// 修改 write offset  对于已经在写了。我们可以认为肯定读了最新的。所以可以也把read_offset 写上
	// 没有会话就创建
	if search1.Id == 0 {
		entry := dao.Conversation{
			UserId:      params.SendId,
			RecvId:      params.RecvId,
			GroupId:     params.GroupId,
			ReadOffset:  params.Seq,
			WriteOffset: params.Seq,
			CreatedTime: time.Now(),
			IsDelete:    0,
		}
		err = tx.Save(&entry).Error
		if err != nil {
			tx.Rollback()
			return err
			// todo
		}
	} else {
		// todo 后续改成 双链更新 排重+有序保证
		// update chat_conversation SET write_offset = {Seq}
		// WHERE  user_id = ? AND target_id = ? AND write_offset = {prevSeq}

		//tx.Raw("UPDATE chat_conversation SET write_offset = ? WHERE user_id = ? AND recv_id = ? AND group_id = ? ",
		//	params.Seq, params.SendId, params.RecvId, params.GroupId)

		err = tx.Model(search1).Updates(map[string]interface{}{"write_offset": params.Seq, "read_offset": params.Seq}).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}

	// 将接受者的会话创建出来
	groupMembers := make([]int, 0, 16)
	tx.Raw("SELECT user_id FROM chat_group_member "+
		"WHERE group_id = ? AND is_delete=0 for update", params.GroupId).Scan(&groupMembers)

	// 查询已经有会话的成员
	convMembers := make([]int, 0, 16)
	tx.Raw("SELECT user_id FROM chat_conversation "+
		"WHERE group_id = ? AND user_id in ? AND is_delete=0 for update", params.GroupId, groupMembers).Scan(&convMembers)

	// 对比两组成员找到出没有会话的，并创建
	var conversations = []*dao.Conversation{}
	for _, gm := range groupMembers {
		isMatch := false
		for _, cm := range convMembers { // 一般群成员不多，直接遍历就行
			if gm == cm {
				isMatch = true
			}
		}
		if isMatch { // 这里包含了自己的item，肯定也会被过滤掉
			continue
		}
		// 没有就添加到创建名单
		entry := dao.Conversation{
			UserId:      gm,
			GroupId:     params.GroupId,
			ReadOffset:  0,
			WriteOffset: 0,
			CreatedTime: time.Now(),
			IsDelete:    0,
		}
		conversations = append(conversations, &entry)
	}
	if len(conversations) > 0 {
		err = tx.CreateInBatches(conversations, len(conversations)).Error
		if err != nil {
			// todo
			log.Println(err)
			tx.Rollback()
			return
		}
	}

	// 判断是否在线
	for _, gm := range groupMembers {
		if gm == params.SendId { // 自己的就不用发送了
			continue
		}
		isOnline := conn_manager.IsAlive(gm)
		onLineMap[gm] = isOnline
		if isOnline {
			c := &dao.Conversation{}
			err = tx.Model(c).Where("user_id = ? AND group_id = ? ", gm, params.GroupId).
				Update("read_offset", params.Seq).Error // 使用最新的 seq 作为 read_offset
			if err != nil {
				log.Fatalln(err)
				tx.Rollback()
				//todo
				return
			}
		}
	}

	tx.Commit()

	// 上面处理没问题就提交偏移量
	kc.Consumer().MarkOffset(msg, "")
	kc.Consumer().CommitOffsets()

	go func() {
		for member := range onLineMap {
			// 发送给websocket
			reply := &conn_handler.SendMsgReply{}
			err = client_service.DefaultConnEngineCli.CallSpecificAddr(
				conn_handler.ConnEngineApiSendMsg,
				conn_handler.SendMsgReq{
					CmdCode:  model.CmdGroup,
					Content:  params.Content,
					Seq:      params.Seq,
					RecvId:   member,
					GroupId:  params.GroupId,
					SendId:   params.SendId,
					SendTime: params.SendTime,
					Type:     params.Type,
				},
				reply,
				conn_manager.Where(member), // 指定 fd 所在的机器提供服务
			)
			if err != nil {
				// todo log
				return
			}
			if reply.Err != nil {
				// todo log
				return
			}
			if reply.Success != true {
				log.Println("send msg to ws fail") // todo log
				return
			}
		}
	}()

	return nil
}
func Private(params *data_handler.SendMsgReq, kc *mkafka.KafkaConsumer, msg *sarama.ConsumerMessage) (err error) {
	isOnline := false

	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return err
	}

	tx = tx.Begin()

	search1 := dao.Conversation{} // for write_offset
	tx.Raw("SELECT * FROM chat_conversation "+
		"WHERE user_id = ? AND recv_id = ? AND group_id = ? for update", params.SendId, params.RecvId, params.GroupId).Scan(&search1)
	// 修改 write offset
	// 没有会话就创建
	if search1.Id == 0 {
		entry := dao.Conversation{
			UserId:  params.SendId,
			RecvId:  params.RecvId,
			GroupId: params.GroupId,
			//ReadOffset:  0,
			WriteOffset: params.Seq,
			CreatedTime: time.Now(),
			IsDelete:    0,
		}
		err = tx.Save(&entry).Error
		if err != nil {
			tx.Rollback()
			return err
			// todo
		}
	} else {
		// todo 后续改成 双链更新 排重+有序保证
		// update chat_conversation SET write_offset = {Seq}
		// WHERE  user_id = ? AND target_id = ? AND write_offset = {prevSeq}

		//tx.Raw("UPDATE chat_conversation SET write_offset = ? WHERE user_id = ? AND recv_id = ? AND group_id = ? ",
		//	params.Seq, params.SendId, params.RecvId, params.GroupId)

		err = tx.Model(search1).Update("write_offset", params.Seq).Error // todo 这里可以加个read_offset,都已经写了。那么可以认为之前的已经读了
		if err != nil {
			tx.Rollback()
			return
		}
	}

	// 将接受者的会话创建出来
	search2 := dao.Conversation{} // for read_offset
	tx.Raw("SELECT * FROM chat_conversation "+
		"WHERE user_id = ? AND recv_id = ? for update", params.RecvId, params.SendId).Scan(&search2)
	// 没有会话就创建
	if search2.Id == 0 {
		entry := dao.Conversation{
			UserId:      params.RecvId,
			RecvId:      params.SendId,
			GroupId:     0,
			ReadOffset:  0,
			WriteOffset: 0,
			CreatedTime: time.Now(),
			IsDelete:    0,
		}
		err = tx.Save(&entry).Error
		if err != nil {
			tx.Rollback()
			return
			// todo
		}
	}

	// 判断是否在线
	isOnline = conn_manager.IsAlive(params.RecvId)

	// 看下目标是否在线，如果在线，需要更新 接受者的read_offset 并发送websocket
	if isOnline {
		c := &dao.Conversation{}
		err = tx.Model(c).Where("user_id = ? AND recv_id = ? AND group_id = ? ", params.RecvId, params.SendId, params.GroupId).Update("read_offset", params.Seq).Error // 使用最新的 seq 作为 read_offset
		if err != nil {
			log.Fatalln(err)
			tx.Rollback()
			//todo
			return
		}
	}
	tx.Commit()

	// 上面处理没问题就提交偏移量
	kc.Consumer().MarkOffset(msg, "")
	kc.Consumer().CommitOffsets()

	if !isOnline {
		return
	}
	// 发送给websocket
	reply := &conn_handler.SendMsgReply{}
	err = client_service.DefaultConnEngineCli.CallSpecificAddr(
		conn_handler.ConnEngineApiSendMsg,
		conn_handler.SendMsgReq{
			CmdCode:  model.CmdPrivate,
			Content:  params.Content,
			Seq:      params.Seq,
			RecvId:   params.RecvId,
			SendId:   params.SendId,
			SendTime: params.SendTime,
			Type:     params.Type,
		},
		reply,
		conn_manager.Where(params.RecvId), // 指定 fd 所在的机器提供服务
	)
	if err != nil {
		// todo
		return err
	}
	if reply.Success != true {
		log.Println("send msg to ws fail") // todo log
	}
	return nil
}
func SendPrivateMsg() {}

func SuccessFunc(kc *mkafka.KafkaConsumer, msg *sarama.ConsumerMessage) (err error) {
	// todo
	log.Println("消费队列处理成功 ", msg)
	return
}
func FailureFunc(kc *mkafka.KafkaConsumer, msg *sarama.ConsumerMessage) (err error) {
	// todo 完善遇到错误时的处理
	log.Println("消费队列处理失败 ", msg)
	return
}
