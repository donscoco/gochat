package controller

import (
	"encoding/json"
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

type GroupController struct{}

func GroupRegister(group *gin.RouterGroup) {
	groupController := &GroupController{}
	group.GET("/list", groupController.List) // 查询用户信息
	group.POST("/create", groupController.Create)
	group.DELETE("/delete/:groupId", groupController.Delete)
	group.PUT("/modify", groupController.Update)
	group.GET("/find/:groupId", groupController.Find)

	group.GET("/members/:groupId", groupController.Get) // 查询群成员用户信息
	group.POST("/invite", groupController.Invite)       // todo 新加入还要在加上conversatino
	group.DELETE("/kick/:groupId", groupController.Kick)
	group.DELETE("/quit/:groupId", groupController.Quit)

}

func (gc *GroupController) List(c *gin.Context) { // domark finish
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	/* DE 逻辑 */
	reply := data_handler.GroupMemberListByUserReply{}
	req := data_handler.GroupMemberListByUserReq{
		UserId: sessionInfo.ID,
	}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberListByUser,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	groupIds := []int{}
	for _, g := range reply.GroupMemberList {
		groupIds = append(groupIds, g.GroupId)
	}

	replyGL := data_handler.GroupListReply{}
	reqGL := data_handler.GroupListReq{GroupIds: groupIds}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupList,
		reqGL,
		&replyGL,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	// 获取到groupList 和groupMember 数据后封装成业务数据model返回
	groupMemberGroupByGroupId := transform.GroupMemberGroupByGroupIdListToMap(reply.GroupMemberList)
	out := transform.TransformGroupList(replyGL.GroupList, groupMemberGroupByGroupId)

	/* DE 逻辑 */
	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}

	//groupMemberSearch := &dao.GroupMember{}
	//
	//groupMemberList, err := groupMemberSearch.ListGroupMemberByUserId(c, tx, sessionInfo.ID)
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}

	//groupIds := []int{}
	//for _, g := range groupMemberList {
	//	groupIds = append(groupIds, g.GroupId)
	//}
	//searchGroup := &dao.Group{}
	//groupList, err := searchGroup.ListGroupByWhere(c, tx, "id in ? AND is_delete= ?", groupIds, 0)
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}

	//groupMemberGroupByGroupId := transform.GroupMemberGroupByGroupIdListToMap(groupMemberList)
	//out := transform.TransformGroupList(groupList, groupMemberGroupByGroupId)
	bl.ResponseSuccess(c, out)
}

func (gc *GroupController) Create(c *gin.Context) { // domark finish

	// 获得请求内容
	params := &model.GroupCreateInput{}
	if err := bl.GetValidParams(c, params, binding.Form); err != nil {
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
	reply := data_handler.GroupCreateReply{}
	req := data_handler.GroupCreateReq{
		UserId:    sessionInfo.ID,
		GroupName: params.GroupName,
		AliasName: sessionInfo.UserName, // 在群中的别名,刚创建默认使用用户名
	}
	err := client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupCreate,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Group == nil || reply.GroupMember == nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	// // 涉及两个表操作 tx.Begin()
	//tx = tx.Begin()
	//
	//insertGroup := &dao.Group{
	//	OwnerId:     sessionInfo.ID,
	//	GroupName:   params.GroupName,
	//	CreatedTime: time.Now(),
	//	IsDelete:    0,
	//}
	//err = insertGroup.Save(c, tx)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//search := &dao.Group{
	//	OwnerId:   sessionInfo.ID,
	//	GroupName: params.GroupName,
	//	IsDelete:  0,
	//}
	//groupEntry, err := search.Find(c, tx, search)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err) // 后面统一一下返回的错误码
	//	return
	//}
	//
	//insertGroupMember := &dao.GroupMember{
	//	GroupId:     groupEntry.Id,
	//	UserId:      sessionInfo.ID,
	//	AliasName:   sessionInfo.UserName, // 第一次创建在群众的别称默认是原来的别称
	//	Remark:      groupEntry.GroupName, // 第一次创建remark 就默认是群名
	//	CreatedTime: time.Now(),
	//	IsDelete:    0,
	//}
	//err = insertGroupMember.Save(c, tx)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err) // 后面统一一下返回的错误码
	//	return
	//}

	// 封装给model 展示
	groupInfo := &model.GroupInfo{
		Id:             reply.Group.Id,
		Name:           reply.Group.GroupName,
		OwnerId:        reply.Group.OwnerId,
		HeadImage:      reply.Group.HeadImage,
		HeadImageThumb: reply.Group.HeadImage,
		Notice:         reply.Group.Notice,
		AliasName:      reply.GroupMember.AliasName,
		Remark:         reply.GroupMember.Remark,
	}

	//out := groupInfo
	bl.ResponseSuccess(c, groupInfo)

}

/*
解散群聊
Request URL: /api/group/delete/25
Request Method: DELETE
{"code":200,"message":"成功","data":null}
*/
func (gc *GroupController) Delete(c *gin.Context) { // domark finish
	groupIdStr := c.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		bl.ResponseError(c, bl.InvalidRequestErrorCode, err)
		return
	}

	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	//tx = tx.Begin()
	//updateGroup := &dao.Group{Id: groupId}
	//err = updateGroup.UpdateById(c, tx, map[string]interface{}{"is_delete": 1})
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err) //  后面统一一下返回的错误码
	//	return
	//}
	//updateGroupMember := &dao.GroupMember{}
	//err = updateGroupMember.UpdateByWhere(c, tx,
	//	map[string]interface{}{"is_delete": 1},
	//	//"user_id= ? AND group_id=?", sessionInfo.ID, groupIdStr,
	//	"group_id=?", groupId,
	//)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err) // 后面统一一下返回的错误码
	//	return
	//}
	//tx.Commit()

	reply := data_handler.GroupDisbandReply{}
	req := data_handler.GroupDisbandReq{GroupId: groupId}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupDisband, //在同一个事务中，要专门提供解散的接口
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
Request URL: /api/group/modify
Request Method: PUT

	{
		"code": 200,
		"message": "成功",
		"data": {
			"id": 21,
			"name": "grouptest",
			"ownerId": 67,
			"headImage": "",
			"headImageThumb": "",
			"notice": "",
			"aliasName": "ironhead",
			"remark": "grouptest100"
		}
	}
*/
func (gc *GroupController) Update(c *gin.Context) { //domark finish
	var err error
	// 获得请求内容
	params := &model.GroupInfo{}
	if err := bl.GetValidParams(c, params, binding.JSON); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	// 获取请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}
	// todo 要检查入参合法性，不能sessionid = ownerid来判断。可能用户搞破坏

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}

	//tx = tx.Begin() //
	//
	//// 群主可以修改的字段
	//if sessionInfo.ID == params.OwnerId {
	//	updateGroup := &dao.Group{Id: params.Id}
	//	mg := map[string]interface{}{
	//		"group_name": params.Name,
	//		"head_image": params.HeadImage,
	//		"notice":     params.Notice,
	//	} // todo 系统这一方需要再做一次校验，看看这个群的owner是否真的是当前用户，这里先简单做，
	//	err = updateGroup.UpdateByWhere(c, tx, mg,
	//		"id = ? AND owner_id = ?", params.Id, sessionInfo.ID)
	//	//err = updateGroup.UpdateById(c, tx, mg)
	//	if err != nil {
	//		tx.Rollback()
	//		bl.ResponseError(c, 2002, err)
	//		return
	//	}
	//}
	//
	//updateGroupMember := &dao.GroupMember{}
	//mgm := map[string]interface{}{
	//	"alias_name": params.AliasName,
	//	"remark":     params.Remark,
	//}
	//err = updateGroupMember.UpdateByWhere(c, tx, mgm,
	//	"user_id=? AND group_id=?", sessionInfo.ID, params.Id)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2003, err) //  后面统一一下返回的错误码
	//	return
	//}
	//tx.Commit()

	//// 群主可以修改的字段
	if sessionInfo.ID == params.OwnerId {
		mg := map[string]interface{}{
			"group_name": params.Name,
			"head_image": params.HeadImage,
			"notice":     params.Notice,
		}
		reply := data_handler.GroupUpdateReply{}
		req := data_handler.GroupUpdateReq{
			GroupId:   params.Id,
			UpdateCol: mg,
		}
		err = client_service.DefaultDataEngineCli.Call(
			data_handler.DataEngineApiGroupUpdate,
			req,
			&reply,
		)
		if err != nil {
			bl.ResponseError(c, bl.RPCErrorCode, err)
			return
		}
	}
	// 群成员可以修改的字段
	mgm := map[string]interface{}{
		"alias_name": params.AliasName,
		"remark":     params.Remark,
	}
	replyGM := data_handler.GroupMemberUpdateReply{}
	reqGM := data_handler.GroupMemberUpdateReq{
		GroupId:   params.Id,
		UserId:    sessionInfo.ID,
		UpdateCol: mgm,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberUpdate,
		reqGM,
		&replyGM,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	groupInfo := model.GroupInfo{
		Id:             params.Id,
		Name:           params.Name,
		OwnerId:        params.OwnerId,
		HeadImage:      params.HeadImage,
		HeadImageThumb: params.HeadImage,
		Notice:         params.Notice,

		AliasName: params.AliasName,
		Remark:    params.Remark,
	}

	bl.ResponseSuccess(c, groupInfo)
}

/*
Request URL: /api/group/find/26
Request Method: GET

	{
		"code": 200,
		"message": "成功",
		"data": {
			"id": 26,
			"name": "grouptest3",
			"ownerId": 68,
			"headImage": "",
			"headImageThumb": "",
			"notice": "",
			"aliasName": "ironhead2",
			"remark": "grouptest3"
		}
	}
*/
func (gc *GroupController) Find(c *gin.Context) { //domark finish
	groupIdStr := c.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
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
	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	//search := dao.Group{}
	//err = tx.Raw("SELECT * FROM chat_group "+
	//	"WHERE id = ?  AND is_delete = 0 ", groupId).Scan(&search).Error
	//if err != nil {
	//}
	//
	//search2 := dao.GroupMember{}
	//err = tx.Raw("SELECT * FROM chat_group_member "+
	//	"WHERE user_id = ? AND group_id = ?  AND is_delete = 0", sessionInfo.ID, groupId).Scan(&search2).Error
	//if err != nil {
	//}

	reply := data_handler.GroupSearchReply{}
	req := data_handler.GroupSearchReq{
		GroupId: groupId,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupSearch,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.GroupInfo == nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	replyGM := data_handler.GroupMemberSearchReply{}
	reqGM := data_handler.GroupMemberSearchReq{
		GroupId: groupId,
		UserId:  sessionInfo.ID,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberSearch,
		reqGM,
		&replyGM,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if replyGM.GroupMemberInfo == nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	groupInfo := model.GroupInfo{
		Id:             reply.GroupInfo.Id,
		Name:           reply.GroupInfo.GroupName,
		OwnerId:        reply.GroupInfo.OwnerId,
		HeadImage:      reply.GroupInfo.HeadImage,
		HeadImageThumb: reply.GroupInfo.HeadImage,
		Notice:         reply.GroupInfo.Notice,
		AliasName:      replyGM.GroupMemberInfo.AliasName,
		Remark:         replyGM.GroupMemberInfo.Remark,
	}

	bl.ResponseSuccess(c, groupInfo)
	return

}

/*
Request URL: /api/group/members/21
Request Method: GET

	{
		"code": 200,
		"message": "成功",
		"data": [{
			"userId": 67,
			"aliasName": "ironhead",
			"headImage": "",
			"quit": false,
			"remark": "grouptest100"
		}, {
			"userId": 68,
			"aliasName": "ironhead2",
			"headImage": "",
			"quit": false,
			"remark": "grouptest2"
		}]
	}
*/
func (gmc *GroupController) Get(c *gin.Context) { // domark finish
	groupIdStr := c.Params.ByName("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {

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
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	//search := &dao.GroupMember{}
	//groupMemberList, err := search.ListGroupMemberByGroupId(c, tx, groupId)
	//if err != nil {
	//
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}

	reply := data_handler.GroupMemberListReply{}
	req := data_handler.GroupMemberListReq{
		GroupId: groupId,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberList,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.GroupMemberList == nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	userIds := make([]int, 0, 16)
	for _, gm := range reply.GroupMemberList {
		userIds = append(userIds, gm.UserId)
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

	out := transform.TransformGroupMemberList(reply.GroupMemberList, userInfoMap)
	bl.ResponseSuccess(c, out)
}

/*
Request URL: http://127.0.0.1:8888/api/group/invite
Request Method: POST

{"groupId":21,"friendIds":[68]}
{"code":200,"message":"成功","data":null}
*/
func (gc *GroupController) Invite(c *gin.Context) { // domark finish

	var err error
	// 获得请求内容
	params := &model.GroupMemberInviteInput{}
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

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}

	// 根据id查询到所有用户的信息，然后插入到groupmember中

	//tx = tx.Begin()
	//searchGroup := &dao.Group{
	//	Id:       params.GroupId,
	//	IsDelete: 0,
	//}
	//groupEntry, err := searchGroup.Find(c, tx, searchGroup)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err) //  后面统一一下返回的错误码
	//	return
	//}
	//
	//search := &dao.User{}
	//userList, err := search.FindByIds(c, tx, params.FriendIds)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	//insertList := make([]dao.GroupMember, 0, 16)
	//for _, u := range userList {
	//	user := dao.GroupMember{
	//		GroupId:     params.GroupId,
	//		UserId:      u.Id,
	//		AliasName:   u.NickName,
	//		Remark:      groupEntry.GroupName,
	//		CreatedTime: time.Now(),
	//		IsDelete:    0,
	//	}
	//	insertList = append(insertList, user)
	//}
	//insert := &dao.GroupMember{}
	//err = insert.CreateInBatches(c, tx, insertList)
	//if err != nil {
	//	tx.Rollback()
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//tx.Commit()

	reply := data_handler.GroupMemberInviteReply{}
	req := data_handler.GroupMemberInviteReq{
		GroupId: params.GroupId,
		UserIds: params.FriendIds,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberInvite,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	bl.ResponseSuccess(c, nil)

}

/*
/api/group/kick/21?userId=68
Request Method: DELETE

{"code":200,"message":"成功","data":null}
*/
func (gmc *GroupController) Kick(c *gin.Context) { // domark finish
	userIdStr := c.Query("userId")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {

		bl.ResponseError(c, 2001, err)
		return
	}
	groupIdStr := c.Params.ByName("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {

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

	// todo 要检查是否是群主，才能t人，这里先默认用户不搞破坏，简单做

	/* DE 逻辑 */
	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	//updateGroupMember := &dao.GroupMember{}
	//mgm := map[string]interface{}{"is_delete": 1}
	//err = updateGroupMember.UpdateByWhere(c, tx, mgm,
	//	"user_id=? AND group_id=?", userId, groupId)
	//if err != nil {
	//	bl.ResponseError(c, 2003, err) //  后面统一一下返回的错误码
	//	return
	//}

	updateCol := map[string]interface{}{
		"is_delete": 1,
	}
	reply := data_handler.GroupMemberUpdateReply{}
	req := data_handler.GroupMemberUpdateReq{
		GroupId: groupId,
		UserId:  userId,

		UpdateCol: updateCol,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberUpdate,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	bl.ResponseSuccess(c, nil)
}

/*
Request URL: /api/group/quit/21
Request Method: DELETE
{"code":200,"message":"成功","data":null}
*/
func (gmc *GroupController) Quit(c *gin.Context) {
	groupIdStr := c.Params.ByName("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {

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
	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//
	//updateGroupMember := &dao.GroupMember{}
	//mgm := map[string]interface{}{"is_delete": 1}
	//err = updateGroupMember.UpdateByWhere(c, tx, mgm,
	//	"user_id=? AND group_id=?", sessionInfo.ID, groupId)
	//if err != nil {
	//	bl.ResponseError(c, 2003, err) //  后面统一一下返回的错误码
	//	return
	//}

	updateCol := map[string]interface{}{
		"is_delete": 1,
	}
	reply := data_handler.GroupMemberUpdateReply{}
	req := data_handler.GroupMemberUpdateReq{
		GroupId: groupId,
		UserId:  sessionInfo.ID,

		UpdateCol: updateCol,
	}
	err = client_service.DefaultDataEngineCli.Call(
		data_handler.DataEngineApiGroupMemberUpdate,
		req,
		&reply,
	)
	if err != nil {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}
	if reply.Code != 200 {
		bl.ResponseError(c, bl.RPCErrorCode, err)
		return
	}

	bl.ResponseSuccess(c, nil)
}
