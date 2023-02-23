package data_manager

import (
	"errors"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/module/data_cache"
	"github.com/donscoco/gochat/internal/module/data_persistence"
	"strconv"
)

func GetContactInfo(uid int, fid int) (contact *dao.Contacts, err error) {

	// 1.查询到数据
	// 2.查询错误
	// 3.查询不到数据
	//	查询数据库，放入缓存
	uidstr := strconv.Itoa(uid)
	fidstr := strconv.Itoa(fid)

	contact, err = data_cache.GetContactsWithCache(uidstr, fidstr)
	if err == nil {
		return contact, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		c, e := data_persistence.LoadContacts(uid, fid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if c.Id > 0 {
			e = data_cache.SetContactsWithCache(uidstr, fidstr, c)
			if e != nil {
				return nil, e
			}
			return c, nil
		}
		// 没有就设置空
		if c.Id == 0 {
			e = data_cache.SetContactsEmptyWithCache(uidstr, fidstr)
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return c, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_CONTACTS_CACHE+uidstr+"_"+fidstr, load)
	if err != nil {
		return nil, err
	}
	contact, ok := result.(*dao.Contacts)
	if !ok {
		err = errors.New("convert error")
	}
	return

}
func SetContactInfo(uid int, fid int, updateCol map[string]interface{}) (err error) {

	// 1.写入db
	// 2.删除缓存，包括list的缓存

	uidstr := strconv.Itoa(uid)
	fidstr := strconv.Itoa(fid)

	err = data_persistence.SaveContacts(uid, fid, updateCol)
	if err != nil {
		return err
	}

	err = data_cache.UnsetContactsWithCache(uidstr, fidstr) //删除
	if err != nil {
		return err
	}

	err = data_cache.UnsetContactsListWithCache(uidstr) //删除
	if err != nil {
		return err
	}

	return err
}
func AddContactInfo(uid int, fid int) (isExist bool, err error) {

	// 1.写入db
	// 2.删除缓存，包括list的缓存

	uidstr := strconv.Itoa(uid)
	fidstr := strconv.Itoa(fid)

	isExist, err = data_persistence.CreateContacts(uid, fid)
	if err != nil {
		return false, err
	}
	if isExist {
		return true, err
	}

	err = data_cache.UnsetContactsWithCache(uidstr, fidstr) //删除
	if err != nil {
		return false, err
	}

	err = data_cache.UnsetContactsListWithCache(uidstr) //删除
	if err != nil {
		return false, err
	}

	return false, err
}

func GetContactList(uid int) (contactList []dao.Contacts, err error) {
	// 好友列表不经常改动，可以直接放string，有好友的修改和添加操作再来删除

	uidstr := strconv.Itoa(uid)
	//contactList = make([]dao.Contacts, 0, 16)
	contactList, err = data_cache.GetContactsListWithCache(uidstr)
	if err == nil {
		return contactList, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		cl, e := data_persistence.LoadContactsList(uid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if len(cl) > 0 {
			e = data_cache.SetContactsListWithCache(uidstr, cl)
			if e != nil {
				return nil, e
			}
			return cl, nil
		}
		// 没有就设置空
		if len(cl) == 0 {
			e = data_cache.SetContactsListEmptyWithCache(uidstr)
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return cl, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_CONTACTS_LIST_CACHE+uidstr, load)
	if err != nil {
		return nil, err
	}
	contactList, ok := result.([]dao.Contacts)
	if !ok {
		err = errors.New("convert error")
	}
	return

}
