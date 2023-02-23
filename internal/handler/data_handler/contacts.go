package data_handler

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/module/data_manager"
)

const (
	DataEngineApiContactsList   = "DataEngine.ContactsList"
	DataEngineApiContactsAdd    = "DataEngine.ContactsAdd"
	DataEngineApiContactsDelete = "DataEngine.ContactsDelete"
	DataEngineApiContactsUpdate = "DataEngine.ContactsUpdate"
	DataEngineApiContactsFind   = "DataEngine.ContactsFind"
)

type ContactsFindReq struct {
	UserId   int
	FriendId int
}
type ContactsFindReply struct {
	Code int
	Msg  string
	Err  error

	ContactInfo *dao.Contacts
}

func (de *DataEngine) ContactsFind(req ContactsFindReq, reply *ContactsFindReply) (err error) {

	// todo log,check param

	contactInfo, err := data_manager.GetContactInfo(req.UserId, req.FriendId)
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
	//contactInfo := dao.Contacts{}
	//err = tx.Raw("SELECT * FROM chat_contacts WHERE user_id = ? AND friend_id = ? AND is_delete = 0 ", req.UserId, req.FriendId).
	//	Scan(&contactInfo).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行错误"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""
	reply.ContactInfo = contactInfo

	return
}

type ContactsUpdateReq struct {
	UserId    int
	FriendId  int
	UpdateCol map[string]interface{}
}
type ContactsUpdateReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) ContactsUpdate(req ContactsUpdateReq, reply *ContactsUpdateReply) (err error) {

	// todo log,check param

	err = data_manager.SetContactInfo(req.UserId, req.FriendId, req.UpdateCol)
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
	//contacts := dao.Contacts{}
	//err = tx.Model(&contacts).Where("user_id = ? AND friend_id = ? ", req.UserId, req.FriendId).
	//	Updates(req.UpdateCol).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行错误"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""

	return
}

type ContactsDeleteReq struct {
	UserId   int
	FriendId int
}
type ContactsDeleteReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) ContactsDelete(req ContactsDeleteReq, reply *ContactsDeleteReply) (err error) {

	// todo log,check param

	m := map[string]interface{}{"is_delete": 1}
	err = data_manager.SetContactInfo(req.UserId, req.FriendId, m)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}

	/* DE 逻辑 */
	////获取连接池的一个连接
	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//contacts := dao.Contacts{}
	//m := map[string]interface{}{"is_delete": 1}
	//err = tx.Model(&contacts).Where("user_id = ? AND friend_id = ? ", req.UserId, req.FriendId).
	//	Updates(m).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行错误"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""

	return
}

type ContactsListReq struct {
	UserId int
}
type ContactsListReply struct {
	Code int
	Msg  string
	Err  error

	ContactsList []dao.Contacts
}

func (de *DataEngine) ContactsList(req ContactsListReq, reply *ContactsListReply) (err error) {

	contactsList, err := data_manager.GetContactList(req.UserId)
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
	//contactsList := make([]dao.Contacts, 0, 16)
	//err = tx.Raw("SELECT * FROM chat_contacts WHERE user_id = ? AND is_delete=0", req.UserId).Scan(&contactsList).Error
	//if err != nil {
	//	reply.Code = 400
	//	reply.Msg = "sql 执行失败"
	//	reply.Err = err
	//	return
	//}

	reply.Code = 200
	reply.Msg = ""
	reply.ContactsList = contactsList

	return
}

type ContactsAddReq struct {
	UserId   int
	FriendId int
}
type ContactsAddReply struct {
	Code int
	Msg  string
	Err  error
}

func (de *DataEngine) ContactsAdd(req ContactsAddReq, reply *ContactsAddReply) (err error) {

	isExist, err := data_manager.AddContactInfo(req.UserId, req.FriendId)
	if err != nil {
		reply.Code = 500
		reply.Msg = "获取DB连接失败"
		reply.Err = err
		return
	}
	if isExist {
		reply.Msg = "已经是好友"
		reply.Code = 200
		return
	}

	//tx, err := gorm.GetGormPool("default")
	//if err != nil {
	//	reply.Code = 500
	//	reply.Msg = "获取DB连接失败"
	//	reply.Err = err
	//	return
	//}
	//
	//tx = tx.Begin()
	//search := dao.Contacts{}
	//err = tx.Raw("SELECT * FROM chat_contacts "+
	//	"WHERE user_id = ? AND friend_id = ? AND is_delete = 0 for update", req.UserId, req.FriendId).Scan(&search).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//if search.FriendId == req.FriendId {
	//	tx.Rollback()
	//	reply.Msg = "已经是好友"
	//	reply.Code = 200
	//	return
	//}
	//
	//// todo 是否开启事务,感觉不用,这里其实只要 拿到 uid，添加查询朋友只是为了拿到nickname，就算中间好友改名nickname不一致也没关系
	//
	//user1 := dao.User{}
	//err = tx.Raw("SELECT * FROM chat_user "+
	//	"WHERE id = ? AND is_delete = 0 ", req.UserId).Scan(&user1).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//user2 := dao.User{}
	//err = tx.Raw("SELECT * FROM chat_user "+
	//	"WHERE id = ? AND is_delete = 0 ", req.FriendId).Scan(&user2).Error
	//if err != nil {
	//	tx.Rollback()
	//	return
	//}
	//
	//insertList := []dao.Contacts{
	//	dao.Contacts{
	//		UserId:         user1.Id,
	//		FriendId:       user2.Id,
	//		FriendNickName: user2.NickName,
	//		Type:           0,
	//		CreatedTime:    time.Now(),
	//		IsDelete:       0,
	//	},
	//	dao.Contacts{
	//		UserId:         user2.Id,
	//		FriendId:       user1.Id,
	//		FriendNickName: user1.NickName,
	//		Type:           0,
	//		CreatedTime:    time.Now(),
	//		IsDelete:       0,
	//	},
	//}
	//err = tx.CreateInBatches(insertList, len(insertList)).Error
	//if err != nil {
	//	return
	//}
	//
	//tx.Commit()

	reply.Code = 200
	reply.Msg = ""

	return
}
