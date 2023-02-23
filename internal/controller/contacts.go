package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/handler/data_handler"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/transform"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"strconv"
)

type ContactsController struct{}

func ContactsRegister(group *gin.RouterGroup) {
	contacts := &ContactsController{}
	group.GET("/list", contacts.List)
	group.POST("/add", contacts.Add)
	group.DELETE("/delete/:uid", contacts.Delete)
	group.PUT("/update", contacts.Update)
	group.GET("/find/:uid", contacts.Find)
}

func (cc *ContactsController) List(c *gin.Context) { // domark finish

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	/* DE 逻辑 */
	reply := data_handler.ContactsListReply{}
	req := data_handler.ContactsListReq{UserId: sessionInfo.ID}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiContactsList,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	userIds := make([]int, 0, 16)
	for _, c := range reply.ContactsList {
		userIds = append(userIds, c.FriendId)
	}

	replyUL := data_handler.UserListReply{}
	reqUL := data_handler.UserListReq{
		UserIds: userIds,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiUserList,
		reqUL,
		&replyUL,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	userInfoMap := make(map[int]dao.User)
	for _, u := range replyUL.UserList {
		userInfoMap[u.Id] = u
	}

	// 封装业务model
	out := transform.TransformContactsList(reply.ContactsList, userInfoMap)
	bl.ResponseSuccess(c, out)
}

/*
添加好友
Request URL: /api/friend/add?friendId=19
Request Method: POST
Content-Type: application/json;charset=UTF-8
{"code":200,"message":"成功","data":null}
*/
func (cc *ContactsController) Add(c *gin.Context) { // domark finish

	// 获得请求内容
	params := &model.ContactsAdd{}
	if err := bl.GetValidParams(c, params, binding.Form); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// 检查
	// 1.检查,写入
	if params.FriendId == sessionInfo.ID {
		bl.ResponseError(c, bl.IllegalOperationErrorCode, errors.New("不能添加自己为好友")) // 不能添加自己为好友
		return
	}

	reply := data_handler.ContactsAddReply{}
	req := data_handler.ContactsAddReq{
		UserId:   sessionInfo.ID,
		FriendId: params.FriendId,
	}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiContactsAdd,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	bl.ResponseSuccess(c, nil)
}

/*
Request URL: /api/friend/delete/19
Request Method: DELETE
{"code":200,"message":"成功","data":null}
*/
func (cc *ContactsController) Delete(c *gin.Context) {

	uidStr := c.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		//
		bl.ResponseError(c, 2001, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	if uid == sessionInfo.ID {
		bl.ResponseMsg(c, bl.IllegalOperationErrorCode, "不能删除自己", nil) // 不能添加自己为好友
		return
	}

	/* DE 逻辑 */
	reply := data_handler.ContactsDeleteReply{}
	req := data_handler.ContactsDeleteReq{
		UserId:   sessionInfo.ID,
		FriendId: uid,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiContactsDelete,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		//
		bl.ResponseError(c, bl.UnknowCode, reply.Err)
		return
	}

	bl.ResponseSuccess(c, nil)
}

func (cc *ContactsController) Update(c *gin.Context) {
	// 获得请求内容

	var err error

	params := &model.ContactsInfo{}
	if err := bl.GetValidParams(c, params, binding.Form); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	reply := data_handler.ContactsUpdateReply{}
	req := data_handler.ContactsUpdateReq{
		UserId:   sessionInfo.ID,
		FriendId: params.Id,
		UpdateCol: map[string]interface{}{
			"friend_nick_name": params.NickName,
		},
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiContactsUpdate,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		//
		bl.ResponseError(c, bl.UnknowCode, reply.Err)
		return
	}

	bl.ResponseSuccess(c, nil)
}

/*
Request URL: http://localhost:8888/api/friend/find/1
Request Method: GET
*/
func (cc *ContactsController) Find(c *gin.Context) {

	uidStr := c.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		//
		bl.ResponseError(c, 2001, err)
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
	reply := data_handler.ContactsFindReply{}
	req := data_handler.ContactsFindReq{
		UserId:   sessionInfo.ID,
		FriendId: uid,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiContactsFind,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseError(c, bl.UnknowCode, reply.Err)
		return
	}

	userInfoReply := data_handler.UserFindReply{}
	userInfoReq := data_handler.UserFindReq{
		UserId: reply.ContactInfo.FriendId,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiUserFind,
		userInfoReq,
		&userInfoReply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseError(c, bl.UnknowCode, reply.Err)
		return
	}

	out := model.ContactsInfo{
		Id:        reply.ContactInfo.FriendId,
		NickName:  reply.ContactInfo.FriendNickName,
		HeadImage: userInfoReply.User.HeadImage,
	}

	bl.ResponseSuccess(c, out)
}
