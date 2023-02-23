package controller

import (
	"encoding/json"
	"fmt"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/handler/data_handler"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/donscoco/gochat/pkg/snowflake"
	"github.com/gin-gonic/gin/binding"
	"strconv"
	"time"

	//"github.com/donscoco/gochat/pkg/gorm"
	//"github.com/donscoco/gochat/pkg/snowflake"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// 根据 uid 和  uid/groupid 和 offset 查找 msg
type MsgController struct{}

func MsgRegister(group *gin.RouterGroup) {
	msgController := &MsgController{}
	group.POST("/private/send", msgController.SendP)
	group.POST("/group/send", msgController.SendG)

	group.POST("/private/pullUnreadMessage", msgController.PullPrivate) // todo
	group.POST("/group/pullUnreadMessage", msgController.PullGroup)     // todo

	group.GET("/private/history", msgController.History) // todo
	group.GET("/group/history", msgController.HistoryG)  // todo

}

/*
发送图片
Request URL: http://localhost/api/message/private/send
Request Method: POST
{
	"id": 0,
	"fileId": 1675888004759,
	"sendId": 68,
	"content": "{\"originUrl\":\"http://localhost/file/box-im/image/20230209/1675888004804.jpg\",\"thumbUrl\":\"http://localhost/file/box-im/image/20230209/1675888004816.jpg\"}",
	"sendTime": 1675888004760,
	"selfSend": true,
	"type": 1,
	"loadStatus": "loading",
	"recvId": 67
}
{"code":200,"message":"成功","data":1061}
*/

/*
发送文件
Request URL: http://localhost/api/message/private/send
Request Method: POST
{
	"id": 0,
	"sendId": 68,
	"content": "{\"name\":\"troubleshooting-kubernetes.pdf\",\"size\":143964,\"url\":\"http://localhost/file/box-im/file/20230209/1675888248836.pdf\"}",
	"sendTime": 1675888248701,
	"selfSend": true,
	"type": 2,
	"loadStatus": "loading",
	"recvId": 67
}
{"code":200,"message":"成功","data":1063}
*/

/*
发送语音
Request URL: http://localhost/api/message/private/send
Request Method: POST
{
	"content": "{\"duration\":3,\"url\":\"http://localhost/file/box-im/file/20230209/1675888437790.wav\"}",
	"type": 3,
	"recvId": 67
}
{"code":200,"message":"成功","data":1063}
*/

/*
发文字
Request URL: http://localhost/api/message/group/send
Request Method: POST
{"content":"test","type":0,"groupId":21}
{"code":200,"message":"成功","data":212}
*/
func (mc *MsgController) SendP(c *gin.Context) {

	params := &model.MsgSendInput{}
	if err := bl.GetValidParams(c, params, binding.JSON); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// todo 排重 或者放到 数据中心去处理

	// 雪花算法生成seq。
	msgSeq := snowflake.GenID()

	// 获得这个会话中的上一个seq，或者说用户在这个会话的offset 然后返回给他
	// 如果不行要不就把offset 放在 cookie 中。？

	// 调用 DE 的rpc，DE 负责写入kafka。
	reply := &data_handler.SendMsgReply{}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiSendMsg,
		data_handler.SendMsgReq{
			Content:  params.Content,
			Seq:      msgSeq,
			RecvId:   params.RecvId,
			SendId:   sessionInfo.ID,
			SendTime: time.Now().Unix(),
			Type:     params.Type,
		},
		reply,
	)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	bl.ResponseSuccess(c, nil)
	return

}

/*
发图片（调用了upload和这个）
Request URL: http://localhost/api/message/group/send
{
  "id": 0,
  "fileId": 1675422358152,
  "sendId": 67,
  "content": "{\"originUrl\":\"http://localhost/file/box-im/image/20230203/1675422358182.jpg\",\"thumbUrl\":\"http://localhost/file/box-im/image/20230203/1675422358192.jpg\"}",
  "sendTime": 1675422358154,
  "selfSend": true,
  "type": 1,
  "loadStatus": "loading",
  "groupId": 21 // recvId:68
}
{"code":200,"message":"成功","data":{"originUrl":"http://localhost/file/box-im/image/20230203/1675422358182.jpg","thumbUrl":"http://localhost/file/box-im/image/20230203/1675422358192.jpg"}}
*/

/*
发语音
Request URL: http://localhost/api/message/group/send
Request Method: POST
{"content":"{\"duration\":4,\"url\":\"http://localhost/file/box-im/file/20230203/1675422465612.wav\"}","type":3,"groupId":21}
{"code":200,"message":"成功","data":214}

私聊
{"content":"{\"duration\":4,\"url\":\"https://donscoco-bucket.oss-cn-guangzhou.aliyuncs.com/user_file/4.wav\"}","type":3,"recvId":1}
*/
func (mc *MsgController) SendG(c *gin.Context) {

	params := &model.MsgSendInput{}
	if err := bl.GetValidParams(c, params, binding.JSON); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// todo 排重

	// 雪花算法生成seq。
	msgSeq := snowflake.GenID()

	// 获得这个会话中的上一个seq，或者说用户在这个会话的offset 然后返回给他
	// 如果不行要不就把offset 放在 cookie 中。？

	// 调用 DE 的rpc，DE 负责写入kafka。

	// 先测试一下写入 DB 和查询DB 和 查询会话
	// 获取连接池的一个连接

	reply := &data_handler.SendMsgReply{}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiSendMsg,
		data_handler.SendMsgReq{
			Content: params.Content,
			Seq:     msgSeq,
			//RecvId:   params.RecvId,
			GroupId:  params.GroupId,
			SendId:   sessionInfo.ID,
			SendTime: time.Now().Unix(),
			Type:     params.Type,
		},
		reply,
	)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	if reply.Success == false {
		bl.ResponseError(c, 2000, err)
		return
	}

	bl.ResponseSuccess(c, nil)
	return

}

/*
Request URL: http://localhost/api/message/private/pullUnreadMessage
Request Method: POST
*/
func (mc *MsgController) PullPrivate(c *gin.Context) {

	// 获取 用户的所有 conversation，查看conversation[] 中的read_offset
	// 将比 read offset 还要大的，recvId 是 当前用户的 消息从mongo中拿出来（限制最多100条）
	// 更新最新的read offset
	// 消息时间排序（小到大）
	// 将消息发送到 对应的channel

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	reply := &data_handler.PullMsgReply{}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiPullMsg,
		data_handler.PullMsgReq{
			UserId: sessionInfo.ID,
		},
		reply,
	)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	bl.ResponseSuccess(c, nil)
	return
}

func (mc *MsgController) PullGroup(c *gin.Context) {
	// 获取用户所在的所有group，查看[group,userId] 的 read_offset
	// 将比 read offset 还要大的，groupId属于用户所在的group集合的消息从mongo中拿出来（多个group）
	// 更新最新的read offset （多个group）
	// 消息时间排序（小到大）
	// 将消息发送到 对应的channel

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	reply := &data_handler.PullMsgReply{}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiPullMsg,
		data_handler.PullMsgReq{
			UserId:  sessionInfo.ID,
			IsGroup: true,
		},
		reply,
	)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	bl.ResponseSuccess(c, nil)
	return
}

/*
Request URL: http://localhost/api/message/private/history?page=1&size=10&friendId=68
Request Method: GET
page=1&size=10&friendId=68

	{
		"code": 200,
		"message": "成功",
		"data": [{
			"id": 961,
			"sendId": 68,
			"recvId": 67,
			"content": "test",
			"type": 0,
			"sendTime": 1675019999000
		}, {
			"id": 960,
			"sendId": 68,
			"recvId": 67,
			"content": "#疑问;",
			"type": 0,
			"sendTime": 1675005172000
		}, {
			"id": 959,
			"sendId": 67,
			"recvId": 68,
			"content": "123456",
			"type": 0,
			"sendTime": 1675005112000
		}]
	}
*/
func (mc *MsgController) History(c *gin.Context) {

	var err error
	pageStr := c.Query("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	sizeStr := c.Query("size")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	fidStr := c.Query("friendId")
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	/* DE 逻辑 */
	reply := &data_handler.HistoryMsgReply{}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiHistory,
		data_handler.HistoryMsgReq{
			Page:     page,
			Size:     size,
			UserId:   int64(sessionInfo.ID),
			FriendId: fid,
		},
		reply,
	)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	bl.ResponseSuccess(c, reply.Records)

}

/*
查看聊天历史
Request URL: http://localhost/api/message/group/history?page=1&size=10&groupId=21
Request Method: GET
page=1&size=10&groupId=21
{"code":200,"message":"成功","data":[{"id":214,"groupId":21,"sendId":67,"content":"{\"duration\":4,\"url\":\"http://localhost/file/box-im/file/20230203/1675422465612.wav\"}","type":3,"sendTime":1675422466000},{"id":213,"groupId":21,"sendId":67,"content":"{\"originUrl\":\"http://localhost/file/box-im/image/20230203/1675422358182.jpg\",\"thumbUrl\":\"http://localhost/file/box-im/image/20230203/1675422358192.jpg\"}","type":1,"sendTime":1675422358000},{"id":212,"groupId":21,"sendId":67,"content":"test","type":0,"sendTime":1675422310000}]}
*/
func (mc *MsgController) HistoryG(c *gin.Context) {

	var err error
	pageStr := c.Query("page")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	sizeStr := c.Query("size")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	gidStr := c.Query("groupId")
	gid, err := strconv.ParseInt(gidStr, 10, 64)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	/* DE 逻辑 */
	reply := &data_handler.HistoryMsgReply{}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiHistory,
		data_handler.HistoryMsgReq{
			Page:   page,
			Size:   size,
			UserId: int64(sessionInfo.ID),
			//FriendId: fid,
			GroupId: gid,
		},
		reply,
	)
	if err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	bl.ResponseSuccess(c, reply.Records)

}
