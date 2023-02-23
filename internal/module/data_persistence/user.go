package data_persistence

import (
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/pkg/gorm"
	"log"
)

func LoadUser(userId int) (userInfo *dao.User, err error) {
	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return nil, err
	}

	userInfo = new(dao.User)
	err = tx.Raw("SELECT * FROM chat_user "+
		"WHERE id = ? AND is_delete = 0 ", userId).Scan(&userInfo).Error
	if err != nil {
		return nil, err
	}

	log.Println("domark DB获取 [user]")

	return userInfo, err
}

func SaveUser(uid int, updateCol map[string]interface{}) (err error) {
	// 获取连接池的一个连接
	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return err
	}

	userInfo := dao.User{}
	err = tx.Model(&userInfo).Where("id = ? AND is_delete = 0", uid).
		Updates(updateCol).Error
	if err != nil {
		return err
	}

	log.Println("domark DB保存 [user]")

	return err
}
