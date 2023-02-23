package data_handler

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/module/data_manager"
)

const (
	DataEngineApiGroupCreate = "DataEngine.GroupCreate" // 创建群组
	DataEngineApiGroupDelete = "DataEngine.GroupDelete" // 删除群组
	DataEngineApiGroupUpdate = "DataEngine.GroupUpdate" // 更新群组
	DataEngineApiGroupSearch = "DataEngine.GroupSearch" // 查询群组

	DataEngineApiGroupMemberCreate = "DataEngine.GroupMemberCreate" // 创建成员
	DataEngineApiGroupMemberDelete = "DataEngine.GroupMemberDelete" // 删除成员
	DataEngineApiGroupMemberUpdate = "DataEngine.GroupMemberUpdate" // 更新成员
	DataEngineApiGroupMemberSearch = "DataEngine.GroupMemberSearch" // 查询成员

	DataEngineApiGroupMemberListByUser = "DataEngine.GroupMemberListByUser" // 根据用户id获取这个用户所属的所有群组
	DataEngineApiGroupList             = "DataEngine.GroupList"             // 根据给定的group数组 查询所有的group
	DataEngineApiGroupMemberList       = "DataEngine.GroupMemberList"       // 根据群组id获取这个群组中的所有用户

	DataEngineApiGroupMemberInvite = "DataEngine.GroupMemberInvite" // 邀请成员
	DataEngineApiGroupMemberKick   = "DataEngine.GroupMemberKick"   // 踢出
	DataEngineApiGroupMemberQuit   = "DataEngine.GroupMemberQuit"   // 退出
	DataEngineApiGroupDisband      = "DataEngine.GroupDisband"      // 解散群组
)

// todo 最后记得检查下 事务中的select 有没有加上 for update

type GroupMemberInviteReq struct {
	GroupId int
	UserIds []int
}
type GroupMemberInviteReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) GroupMemberInvite(req GroupMemberInviteReq, reply *GroupMemberInviteReply) (err error) {

	err = data_manager.GroupMemberInvite(req.GroupId, req.UserIds)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	////todo 是否不用开启事务？ 前面的两句查询只是为了拿到groupName 和 userName，用于初始化默认的群备注 这种就算不一致也没关系。
	//
	//// 查询群组
	//// 查询所有用户
	//// 创建所有groupmember
	//
	//searchGroup := dao.Group{}
	//err = tx.Raw("SELECT * FROM chat_group WHERE id = ? AND is_delete =0 ", req.GroupId).Scan(&searchGroup).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//userList := make([]dao.User, 0, 16)
	//err = tx.Raw("SELECT * FROM chat_user WHERE id in ? AND is_delete =0 ", req.UserIds).Scan(&userList).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//groupMemberList := make([]dao.GroupMember, 0, 16)
	//for _, u := range userList {
	//	user := dao.GroupMember{
	//		GroupId:     req.GroupId,
	//		UserId:      u.Id,
	//		AliasName:   u.UserName,
	//		Remark:      searchGroup.GroupName,
	//		CreatedTime: time.Now(),
	//		IsDelete:    0,
	//	}
	//	groupMemberList = append(groupMemberList, user)
	//}
	//
	//err = tx.CreateInBatches(groupMemberList, len(groupMemberList)).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""

	return
}

type GroupMemberDeleteReq struct {
	GroupId int
	UserId  int
}
type GroupMemberDeleteReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) GroupMemberDelete(req GroupMemberDeleteReq, reply *GroupMemberDeleteReply) (err error) {

	m := map[string]interface{}{"is_delete": 1}
	err = data_manager.SetGroupMember(req.GroupId, req.UserId, m)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//updateGroupMember := dao.GroupMember{}
	//m := map[string]interface{}{"is_delete": 1}
	//err = tx.Model(&updateGroupMember).Where("group_id = ? AND user_id = ?", req.GroupId, req.UserId).Updates(m).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""

	return
}

type GroupMemberUpdateReq struct {
	GroupId int
	UserId  int

	UpdateCol map[string]interface{}
}
type GroupMemberUpdateReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) GroupMemberUpdate(req GroupMemberUpdateReq, reply *GroupMemberUpdateReply) (err error) {

	err = data_manager.SetGroupMember(req.GroupId, req.UserId, req.UpdateCol)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//updateGroupMember := dao.GroupMember{}
	//err = tx.Model(&updateGroupMember).Where("group_id = ? AND user_id = ?", req.GroupId, req.UserId).Updates(req.UpdateCol).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""

	return
}

type GroupMemberSearchReq struct {
	GroupId int
	UserId  int
}
type GroupMemberSearchReply struct {
	Code int
	Msg  string
	Err  error

	GroupMemberInfo *dao.GroupMember
}

func (de *DataEngine) GroupMemberSearch(req GroupMemberSearchReq, reply *GroupMemberSearchReply) (err error) {

	groupMember, err := data_manager.GetGroupMember(req.GroupId, req.UserId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//searchGroupMember := dao.GroupMember{}
	//err = tx.Raw("SELECT * FROM chat_group_member WHERE group_id = ? AND user_id = ? AND is_delete = 0 ", req.GroupId, req.UserId).Scan(&searchGroupMember).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""
	reply.GroupMemberInfo = groupMember

	return
}

///

type GroupCreateReq struct {
	UserId    int
	GroupName string

	AliasName string // 在群中的别名
}
type GroupCreateReply struct {
	Code int
	Msg  string
	Err  error

	Group       *dao.Group
	GroupMember *dao.GroupMember
}

func (de *DataEngine) GroupCreate(req GroupCreateReq, reply *GroupCreateReply) (err error) {

	group, groupMember, err := data_manager.AddGroup(req.UserId, req.AliasName, req.GroupName)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//insertGroup := dao.Group{
	//	OwnerId:     req.UserId,
	//	GroupName:   req.GroupName,
	//	CreatedTime: time.Now(),
	//	IsDelete:    0,
	//}
	//err = tx.Save(&insertGroup).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//err = tx.Raw("SELECT * FROM chat_group "+
	//	"WHERE owner_id = ? AND group_name = ? AND is_delete = 0 ", req.UserId, req.GroupName).Scan(&insertGroup).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//insertGroupMember := dao.GroupMember{
	//	GroupId:     insertGroup.Id,
	//	UserId:      req.UserId,
	//	AliasName:   req.AliasName,         // 第一次创建在群众的别称默认是原来的别称
	//	Remark:      insertGroup.GroupName, // 第一次创建remark 就默认是群名
	//	CreatedTime: time.Now(),
	//	IsDelete:    0,
	//}
	//err = tx.Save(&insertGroupMember).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""
	reply.Group = group
	reply.GroupMember = groupMember

	return
}

type GroupDeleteReq struct {
	GroupId int
}
type GroupDeleteReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) GroupDelete(req GroupDeleteReq, reply *GroupDeleteReply) (err error) {

	m := map[string]interface{}{"is_delete": 1}
	err = data_manager.SetGroupInfo(req.GroupId, m)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//m := map[string]interface{}{"is_delete": 1}
	//updateGroup := dao.Group{}
	//err = tx.Model(&updateGroup).Where("id = ?", req.GroupId).Updates(m).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//updateGroupMember := dao.GroupMember{}
	//err = tx.Model(&updateGroupMember).Where("group_id = ? ", req.GroupId).Updates(m).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""

	return
}

type GroupUpdateReq struct {
	GroupId int

	UpdateCol map[string]interface{}
}
type GroupUpdateReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) GroupUpdate(req GroupUpdateReq, reply *GroupUpdateReply) (err error) {

	err = data_manager.SetGroupInfo(req.GroupId, req.UpdateCol)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//updateGroup := dao.Group{}
	//err = tx.Model(&updateGroup).Where("id = ?", req.GroupId).Updates(req.UpdateCol).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""

	return
}

type GroupSearchReq struct {
	GroupId int
}
type GroupSearchReply struct {
	Code int
	Msg  string
	Err  error

	GroupInfo *dao.Group
}

func (de *DataEngine) GroupSearch(req GroupSearchReq, reply *GroupSearchReply) (err error) {

	groupInfo, err := data_manager.GetGroupInfo(req.GroupId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//searchGroup := dao.Group{}
	//err = tx.Raw("SELECT * FROM chat_group WHERE id = ? AND is_delete = 0 ", req.GroupId).Scan(&searchGroup).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""
	reply.GroupInfo = groupInfo

	return
}

////

type GroupMemberListByUserReq struct {
	UserId int
}
type GroupMemberListByUserReply struct {
	Code int
	Msg  string
	Err  error

	GroupMemberList []dao.GroupMember
}

func (de *DataEngine) GroupMemberListByUser(req GroupMemberListByUserReq, reply *GroupMemberListByUserReply) (err error) {

	groupMemberInfoList, err := data_manager.GetGroupListByUser(req.UserId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//groupMemberInfoList := make([]dao.GroupMember, 0, 16)
	//err = tx.Raw("SELECT * FROM chat_group_member WHERE user_id = ? AND is_delete = 0 ", req.UserId).
	//	Scan(&groupMemberInfoList).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行错误"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""
	reply.GroupMemberList = groupMemberInfoList

	return
}

type GroupMemberListReq struct {
	GroupId int
}
type GroupMemberListReply struct {
	Code int
	Msg  string
	Err  error

	GroupMemberList []dao.GroupMember
}

func (de *DataEngine) GroupMemberList(req GroupMemberListReq, reply *GroupMemberListReply) (err error) {

	groupMemberInfoList, err := data_manager.GetGroupMemberList(req.GroupId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//groupMemberInfoList := make([]dao.GroupMember, 0, 16)
	//err = tx.Raw("SELECT * FROM chat_group_member WHERE group_id = ? AND is_delete = 0 ", req.GroupId).
	//	Scan(&groupMemberInfoList).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行错误"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""
	reply.GroupMemberList = groupMemberInfoList

	return
}

type GroupListReq struct {
	GroupIds []int
}
type GroupListReply struct {
	Code int
	Msg  string
	Err  error

	GroupList []dao.Group
}

func (de *DataEngine) GroupList(req GroupListReq, reply *GroupListReply) (err error) {

	groupInfoList, err := data_manager.GetGroupList(req.GroupIds)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	///* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//groupInfoList := make([]dao.Group, 0, 16)
	//err = tx.Raw("SELECT * FROM chat_group WHERE id in ? AND is_delete = 0 ", req.GroupIds).
	//	Scan(&groupInfoList).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行错误"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""
	reply.GroupList = groupInfoList

	return
}

// GroupDisband
type GroupDisbandReq struct {
	GroupId int
}
type GroupDisbandReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) GroupDisband(req GroupDisbandReq, reply *GroupDisbandReply) (err error) {

	err = data_manager.DisbandGroup(req.GroupId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	/* DE 逻辑 */
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//
	//m := map[string]interface{}{"is_delete": 1}
	//updateGroup := dao.Group{}
	//err = tx.Model(&updateGroup).Where("id = ?", req.GroupId).Updates(m).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//updateGroupMember := dao.GroupMember{}
	//err = tx.Model(&updateGroupMember).Where("group_id = ? ", req.GroupId).Updates(m).Error
	//if err != nil {
	//	tx.Rollback()
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""

	return
}
