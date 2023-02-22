package data_handler

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/module/data_manager"
	"github.com/donscoco/gochat/pkg/gorm"
)

const (
	DataEngineApiUserFindByNickName = "DataEngine.UserFindByNickName"
	DataEngineApiUserInfo           = "DataEngine.UserInfo"
	DataEngineApiUserFind           = "DataEngine.UserFind"
	DataEngineApiUserUpdate         = "DataEngine.UserUpdate"

	DataEngineApiUserList = "DataEngine.UserList" // 给定用户id数组，查询所有用户信息
)

type UserListReq struct {
	UserIds []int
}
type UserListReply struct {
	Code int
	Msg  string
	Err  error

	UserList []dao.User
}

func (de *DataEngine) UserList(req UserListReq, reply *UserListReply) (err error) {
	/* DE 逻辑 */

	// todo 查询缓存
	// 没有缓存就查询数据库，并写入缓存

	userList, err := data_manager.GetUserList(req.UserIds)
	if err != nil {
		reply.Code = 500
		reply.Msg = "查询数据管理失败"
		reply.Err = err
		// todo log
		return
	}

	//// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	// todo
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//userList := make([]dao.User, 0, 16)
	//err = tx.Raw("SELECT * FROM chat_user "+
	//	"WHERE id in ? AND is_delete = 0 ", req.UserIds).Scan(&userList).Error
	//if err != nil {
	//	// todo
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}

	reply.UserList = userList
	reply.Code = 200
	return nil

}

type UserFindReq struct {
	UserId int
}
type UserFindReply struct {
	Code int
	Msg  string
	Err  error

	User *dao.User
}

func (de *DataEngine) UserFind(req UserFindReq, reply *UserFindReply) (err error) {
	/* DE 逻辑 */
	userInfo, err := data_manager.GetUserInfo(req.UserId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "查询数据管理失败"
		reply.Err = err
		// todo log
		return
	}

	reply.User = userInfo
	reply.Code = 200
	return nil

	/////////////////////////////////////////////////////

	// todo 查询缓存
	// 没有缓存就查询数据库，并写入缓存

	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	// todo
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//userInfo := dao.User{}
	//err = tx.Raw("SELECT * FROM chat_user "+
	//	"WHERE id = ? AND is_delete = 0 ", req.UserId).Scan(&userInfo).Error
	//if err != nil {
	//	// todo
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//reply.User = &userInfo
	//reply.Code = 200
	//return nil

}

type UserUpdateReq struct {
	UserId    int
	UpdateCol map[string]interface{}
}
type UserUpdateReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) UserUpdate(req UserUpdateReq, reply *UserUpdateReply) (err error) {
	/* DE 逻辑 */

	// todo 参数检查

	err = data_manager.SetUserInfo(req.UserId, req.UpdateCol)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
	}

	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	// todo
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//userInfo := dao.User{}
	//err = tx.Model(&userInfo).Where("id = ? AND is_delete = 0", req.UserId).
	//	Updates(req.UpdateCol).Error
	//if err != nil {
	//	// todo
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200

	return nil

}

type UserInfoReq struct {
	UserId int
}
type UserInfoReply struct {
	Code int
	Msg  string
	Err  error

	User *dao.User
}

func (de *DataEngine) UserInfo(req UserInfoReq, reply *UserInfoReply) (err error) {
	/* DE 逻辑 */

	userInfo, err := data_manager.GetUserInfo(req.UserId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "查询数据管理失败"
		reply.Err = err
		// todo log
		return
	}

	reply.User = userInfo
	reply.Code = 200
	return nil

	/////////////////////////////////////////////////////

	// todo 查询缓存

	// 获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	// todo
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//userInfo := dao.User{}
	//err = tx.Raw("SELECT * FROM chat_user "+
	//	"WHERE id = ? AND is_delete = 0", req.UserId).Scan(&userInfo).Error
	//if err != nil {
	//	// todo
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}
	//
	//reply.Code = 200
	//reply.User = &userInfo
	//
	//return nil

}

type FindByNickNameReq struct {
	NickName string
}
type FindByNickNameReply struct {
	Code int
	Msg  string
	Err  error

	UserList []dao.User
}

func (de *DataEngine) UserFindByNickName(req FindByNickNameReq, reply *FindByNickNameReply) (err error) {
	/* DE 逻辑 */

	// todo 查询缓存

	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		// todo
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	userList := make([]dao.User, 0, 16)
	err = tx.Raw("SELECT * FROM chat_user "+
		"WHERE nick_name like ? AND is_delete = 0", req.NickName+"%").
		Scan(&userList).Error
	if err != nil {
		// todo
		reply.Code = 400
		reply.Msg = "sql 执行失败"
		reply.Err = err
		return
	}

	reply.Code = 200
	reply.UserList = userList
	return nil
}
