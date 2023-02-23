package data_persistence

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/pkg/gorm"
	"log"
	"time"
)

func LoadGroup(groupId int) (groupInfo *dao.Group, err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	tx = tx.Begin()

	groupInfo = new(dao.Group)
	err = tx.Raw("SELECT * FROM chat_group WHERE id = ? AND is_delete = 0 ", groupId).Scan(groupInfo).Error
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	log.Println("domark DB获取 [group]")

	return
}
func SaveGroup(groupId int, updateCol map[string]interface{}) (err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	tx = tx.Begin()

	updateGroup := dao.Group{}
	err = tx.Model(&updateGroup).Where("id = ?", groupId).Updates(updateCol).Error
	if err != nil {
		tx.Rollback()
		//reply.Code = 400
		//reply.Msg = "sql 执行失败"
		//reply.Err = err
		return
	}

	tx.Commit()

	log.Println("domark DB保存 [group]")

	return
}
func CreateGroup(uid int, aliasName, groupName string) (insertGroup *dao.Group, insertGroupMember *dao.GroupMember, err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return nil, nil, err
	}

	tx = tx.Begin()

	insertGroup = &dao.Group{
		OwnerId:     uid,
		GroupName:   groupName,
		CreatedTime: time.Now(),
		IsDelete:    0,
	}
	err = tx.Save(insertGroup).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	err = tx.Raw("SELECT * FROM chat_group "+
		"WHERE owner_id = ? AND group_name = ? AND is_delete = 0 ", uid, groupName).Scan(&insertGroup).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	insertGroupMember = &dao.GroupMember{
		GroupId:     insertGroup.Id,
		UserId:      uid,
		AliasName:   aliasName,             // 第一次创建在群众的别称默认是原来的别称
		Remark:      insertGroup.GroupName, // 第一次创建remark 就默认是群名
		CreatedTime: time.Now(),
		IsDelete:    0,
	}
	err = tx.Save(insertGroupMember).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tx.Commit()

	log.Println("domark DB创建 [group]")

	return nil, nil, err
}
func DistoryGroup(groupId int) (err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		//reply.Code = 500
		//reply.Msg = "获取DB连接失败"
		//reply.Err = err
		return
	}

	tx = tx.Begin()

	m := map[string]interface{}{"is_delete": 1}

	updateGroup := dao.Group{}
	err = tx.Model(&updateGroup).Where("id = ?", groupId).Updates(m).Error
	if err != nil {
		tx.Rollback()
		//reply.Code = 400
		//reply.Msg = "sql 执行失败"
		//reply.Err = err
		return
	}
	updateGroupMember := dao.GroupMember{}
	err = tx.Model(&updateGroupMember).Where("group_id = ? ", groupId).Updates(m).Error
	if err != nil {
		tx.Rollback()
		//reply.Code = 400
		//reply.Msg = "sql 执行失败"
		//reply.Err = err
		return
	}

	tx.Commit()

	log.Println("domark DB解散 [group]")

	return err
}

func LoadGroupMember(uid, gid int) (searchGroupMember *dao.GroupMember, err error) {

	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	searchGroupMember = new(dao.GroupMember)
	err = tx.Raw("SELECT * FROM chat_group_member WHERE group_id = ? AND user_id = ? AND is_delete = 0 ", gid, uid).
		Scan(searchGroupMember).Error
	if err != nil {
		tx.Rollback()
		//reply.Code = 400
		//reply.Msg = "sql 执行失败"
		//reply.Err = err
		return
	}

	log.Println("domark DB获取 [groupMember]")

	return searchGroupMember, nil

}
func SaveGroupMember(gid, uid int, updateCol map[string]interface{}) (err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	tx = tx.Begin()

	updateGroupMember := dao.GroupMember{}
	err = tx.Model(&updateGroupMember).Where("group_id = ? AND user_id = ?", gid, uid).Updates(updateCol).Error
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()

	log.Println("domark DB保存 [groupMember]")

	return
}
func CreateGroupMemberInBatch(members []dao.GroupMember) (err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	err = tx.CreateInBatches(members, len(members)).Error
	if err != nil {
		return
	}

	log.Println("domark DB创建 [groupMember]")

	return
}

func LoadGroupMemberList(groupId int) (members []dao.GroupMember, err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	members = make([]dao.GroupMember, 0, 16)
	err = tx.Raw("SELECT * FROM chat_group_member WHERE group_id = ? AND is_delete = 0 ", groupId).
		Scan(&members).Error
	if err != nil {
		return
	}

	log.Println("domark DB获取 [groupMember list]")

	return
}
func LoadGroupListByUId(uid int) (groups []dao.GroupMember, err error) {
	/* DE 逻辑 */
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		//reply.Code = 500
		//reply.Msg = "获取DB连接失败"
		//reply.Err = err
		return
	}

	groups = make([]dao.GroupMember, 0, 16)
	err = tx.Raw("SELECT * FROM chat_group_member WHERE user_id = ? AND is_delete = 0 ", uid).
		Scan(&groups).Error
	if err != nil {
		return nil, err
	}

	log.Println("domark DB获取 [group list]")

	return
}
