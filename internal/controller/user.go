package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/handler/data_handler"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/conn_manager"
	"github.com/donscoco/gochat/internal/module/transform"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type UserController struct{}

func UserRegister(group *gin.RouterGroup) {
	userController := &UserController{}
	group.GET("/self", userController.SelfInfo) // 查询用户信息
	group.GET("/findByNickName", userController.FindByNickName)
	group.PUT("/update", userController.Update)
	group.GET("/find/:uid", userController.Find)
	group.GET("/online", userController.Online)
}

func (u *UserController) SelfInfo(c *gin.Context) { // domark finish

	// 准备参数
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	userSessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), userSessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	/* DE 数据 */
	reply := data_handler.UserInfoReply{}
	req := data_handler.UserInfoReq{UserId: userSessionInfo.ID}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiUserInfo,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseMsg(c, bl.UnknowCode, "数据服务未启动", nil)
		return
	}
	//if reply.User == nil {
	//	bl.ResponseError(c, bl.NilCode, err)
	//	return
	//}

	// 查询获得数据后进行封装，// 转化业务数据model
	out := &model.UserInfo{
		Id:             reply.User.Id,
		UserName:       reply.User.UserName,
		NickName:       reply.User.NickName,
		Sex:            reply.User.Sex,
		Signature:      reply.User.Signature,
		HeadImage:      reply.User.HeadImage,
		HeadImageThumb: reply.User.HeadImage, // todo 后续考虑添加缩略图
		Online:         conn_manager.IsAlive(reply.User.Id),
	}
	bl.ResponseSuccess(c, out)
}

/*
/api/user/findByNickName?nickName=
GET
*/
func (uc *UserController) FindByNickName(c *gin.Context) { // domark finish

	// 准备参数
	params := &model.FindUserInput{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, bl.InvalidRequestErrorCode, err)
		return
	}

	/* DE 数据 */
	reply := data_handler.FindByNickNameReply{}
	req := data_handler.FindByNickNameReq{NickName: params.NickName}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiUserFindByNickName,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	// 转化业务数据model
	out := transform.TransformUserList(reply.UserList)
	bl.ResponseSuccess(c, out)
}

// 修改用户信息
/*
Request URL: /api/user/update
Request Method: PUT

{"id":67,"userName":"ironhead","nickName":"ironhead","sex":0,"signature":"test","headImage":"http://localhost/file/box-im/image/20230203/1675414276279.png","headImageThumb":"http://localhost/file/box-im/image/20230203/1675414276293.png","online":null}
*/
func (u *UserController) Update(c *gin.Context) { // domark finish

	// 准备参数
	params := &model.UserInfo{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	userSessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), userSessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// todo 要做一些安全检查
	if params.Id != userSessionInfo.ID {
		bl.ResponseError(c, bl.IllegalOperationErrorCode, errors.New("不能修改别人的信息"))
		return
	}

	/* DE 数据 */
	reply := data_handler.UserUpdateReply{}
	updateCol := map[string]interface{}{
		"user_name":  params.UserName,
		"nick_name":  params.NickName,
		"sex":        params.Sex,
		"signature":  params.Signature,
		"head_image": params.HeadImage,
	}
	req := data_handler.UserUpdateReq{
		UserId:    params.Id,
		UpdateCol: updateCol,
	}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiUserUpdate,
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
Request URL: /api/user/find/68
Request Method: GET
*/
func (uc *UserController) Find(c *gin.Context) { // domark finish

	// 准备参数
	uidStr := c.Param("uid")
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		// todo
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

	/* DE 数据 */
	reply := data_handler.UserFindReply{}
	req := data_handler.UserFindReq{UserId: uid}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiUserFind,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseMsg(c, bl.UnknowCode, "数据服务未启动", nil)
		return
	}

	// 封装一下业务逻辑
	out := model.UserInfo{
		Id:             reply.User.Id,
		UserName:       reply.User.UserName,
		NickName:       reply.User.NickName,
		Sex:            reply.User.Sex,
		Signature:      reply.User.Signature,
		HeadImage:      reply.User.HeadImage,
		HeadImageThumb: reply.User.HeadImage,
		Online:         conn_manager.IsAlive(reply.User.Id),
	}
	bl.ResponseSuccess(c, out)
}

func (uc *UserController) Online(c *gin.Context) { // domark finish
	userIds := c.Query("userIds")
	uidStrs := strings.Split(userIds, ",")

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// 查询在线状态
	onlineMember := make([]int, 0, len(uidStrs))
	for _, uidstr := range uidStrs {
		uid, err := strconv.Atoi(uidstr)
		if err != nil {
			continue
		}
		if conn_manager.IsAlive(uid) {
			onlineMember = append(onlineMember, uid)
		}
	}

	onlineMember = append(onlineMember, sessionInfo.ID)

	if len(onlineMember) > 0 {
		bl.ResponseSuccess(c, onlineMember)
		return
	}
	bl.ResponseMsg(c, bl.OffLineCode, "未登录", nil)
}
