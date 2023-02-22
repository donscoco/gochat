package data_manager

import (
	"errors"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/module/data_cache"
	"github.com/donscoco/gochat/internal/module/data_persistence"
	"strconv"
)

/*
负责做DB数据和缓存数据之间的数据结构转化管理等

查询cache
sigleflight
查询DB
*/

func GetUserInfo(uid int) (user *dao.User, err error) {

	// 1.查询到数据
	// 2.查询错误
	// 3.查询不到数据
	//	查询数据库，放入缓存

	user, err = data_cache.GetUserWithCache(strconv.Itoa(uid))
	if err == nil {
		return user, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		u, e := data_persistence.LoadUser(uid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if u.Id > 0 {
			e = data_cache.SetUserWithCache(strconv.Itoa(u.Id), u)
			if e != nil {
				return nil, e
			}
			return u, nil
		}
		// 没有就设置空
		if u.Id == 0 {
			e = data_cache.SetUserEmptyWithCache(strconv.Itoa(uid))
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return u, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_USER_CACHE+strconv.Itoa(uid), load)
	if err != nil {
		return nil, err
	}
	user, ok := result.(*dao.User)
	if !ok {
		err = errors.New("convert error")
	}
	return

}

func SetUserInfo(uid int, updateCol map[string]interface{}) (err error) {

	// 1.写入db
	// 2.删除缓存

	err = data_persistence.SaveUser(uid, updateCol)
	if err != nil {
		return err
	}

	err = data_cache.UnsetUserWithCache(strconv.Itoa(uid)) //删除
	if err != nil {
		return err
	}

	return err
}

func GetUserList(uids []int) (userList []dao.User, err error) {

	userList = make([]dao.User, 0, 16)
	for _, uid := range uids {
		u, err := GetUserInfo(uid)
		if err != nil {
			continue
		}
		if u.Id == 0 { // 数据为空的情况
			continue
		}
		userList = append(userList, *u)
	}
	return userList, err
}
