package data_persistence

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/pkg/gorm"
	"log"
	"time"
)

func LoadContacts(userId int, friendId int) (contactInfo *dao.Contacts, err error) {
	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	contactInfo = new(dao.Contacts)
	err = tx.Raw("SELECT * FROM chat_contacts WHERE user_id = ? AND friend_id = ? AND is_delete = 0 ", userId, friendId).
		Scan(contactInfo).Error
	if err != nil {
		return nil, err
	}

	log.Println("domark DB获取 [contacts]")

	return contactInfo, err
}

func SaveContacts(uid int, fid int, updateCol map[string]interface{}) (err error) {
	/* DE 逻辑 */
	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}
	contacts := dao.Contacts{}
	err = tx.Model(&contacts).Where("user_id = ? AND friend_id = ? ", uid, fid).
		Updates(updateCol).Error
	if err != nil {
		// todo
		return
	}

	log.Println("domark DB保存 [contacts]")

	return err
}

func CreateContacts(uid int, fid int) (exist bool, err error) {
	/* DE 逻辑 */
	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	tx = tx.Begin()
	search := dao.Contacts{}
	err = tx.Raw("SELECT * FROM chat_contacts "+
		"WHERE user_id = ? AND friend_id = ? AND is_delete = 0 for update", uid, fid).Scan(&search).Error
	if err != nil {
		// todo
		tx.Rollback()
		return
	}
	if search.FriendId == fid {
		tx.Rollback()
		// 已经是好友了
		return true, err
	}

	user1 := dao.User{}
	err = tx.Raw("SELECT * FROM chat_user "+
		"WHERE id = ? AND is_delete = 0 ", uid).Scan(&user1).Error
	if err != nil {
		// todo
		tx.Rollback()
		return
	}
	user2 := dao.User{}
	err = tx.Raw("SELECT * FROM chat_user "+
		"WHERE id = ? AND is_delete = 0 ", fid).Scan(&user2).Error
	if err != nil {
		//todo
		tx.Rollback()
		return
	}

	insertList := []dao.Contacts{
		dao.Contacts{
			UserId:         user1.Id,
			FriendId:       user2.Id,
			FriendNickName: user2.NickName,
			Type:           0,
			CreatedTime:    time.Now(),
			IsDelete:       0,
		},
		dao.Contacts{
			UserId:         user2.Id,
			FriendId:       user1.Id,
			FriendNickName: user1.NickName,
			Type:           0,
			CreatedTime:    time.Now(),
			IsDelete:       0,
		},
	}
	err = tx.CreateInBatches(insertList, len(insertList)).Error
	if err != nil {
		tx.Rollback()
		// todo
		return
	}

	tx.Commit()

	log.Println("domark DB创建 [contacts]")

	return false, err
}

// ext

func LoadContactsList(userId int) (contactsList []dao.Contacts, err error) {
	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	contactsList = make([]dao.Contacts, 0, 16)
	err = tx.Raw("SELECT * FROM chat_contacts WHERE user_id = ? AND is_delete=0", userId).Scan(&contactsList).Error
	if err != nil {
		return
	}

	log.Println("domark DB获取 [contactsList]")

	return contactsList, err
}
