package data_manager

import (
	"errors"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/module/data_cache"
	"github.com/donscoco/gochat/internal/module/data_persistence"
	"strconv"
	"time"
)

func GetGroupInfo(gid int) (group *dao.Group, err error) {
	// 查询

	// 1.查询到数据
	// 2.查询错误
	// 3.查询不到数据
	//	查询数据库，放入缓存

	gidstr := strconv.Itoa(gid)

	group, err = data_cache.GetGroupWithCache(gidstr)
	if err == nil {
		return group, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		g, e := data_persistence.LoadGroup(gid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if g.Id > 0 {
			e = data_cache.SetGroupWithCache(gidstr, g)
			if e != nil {
				return nil, e
			}
			return g, nil
		}
		// 没有就设置空
		if g.Id == 0 {
			e = data_cache.SetGroupEmptyWithCache(gidstr)
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return g, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_GROUP_CACHE+gidstr, load)
	if err != nil {
		return nil, err
	}
	group, ok := result.(*dao.Group)
	if !ok {
		err = errors.New("convert error")
	}

	return
}

func SetGroupInfo(gid int, updateCol map[string]interface{}) (err error) {
	// 1.写入db
	// 2.删除缓存

	gidstr := strconv.Itoa(gid)

	err = data_persistence.SaveGroup(gid, updateCol)
	if err != nil {
		return err
	}

	err = data_cache.UnsetGroupWithCache(gidstr) //删除
	if err != nil {
		return err
	}

	gml, err := data_cache.GetGroupMemberListWithCache(gidstr)
	if err != nil {
		return err
	}
	for _, gm := range gml {
		data_cache.UnsetGroupListWithCache(strconv.Itoa(gm.UserId))
	}

	err = data_cache.UnsetGroupMemberListWithCache(gidstr) //删除
	if err != nil {
		return err
	}

	return
}

/* group member */
func GetGroupMember(gid int, uid int) (groupmember *dao.GroupMember, err error) {
	// 查询缓存
	// 正常错误
	// 找不到
	// sigleflight 查询db，放入缓存

	gidstr := strconv.Itoa(gid)
	uidstr := strconv.Itoa(uid)

	groupmember, err = data_cache.GetGroupMemberWithCache(uidstr, gidstr)
	if err == nil {
		return groupmember, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		gm, e := data_persistence.LoadGroupMember(uid, gid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if gm.Id > 0 {
			e = data_cache.SetGroupMemberWithCache(uidstr, gidstr, gm)
			if e != nil {
				return nil, e
			}
			return gm, nil
		}
		// 没有就设置空
		if gm.Id == 0 {
			e = data_cache.SetGroupMemberEmptyWithCache(uidstr, gidstr)
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return gm, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_GROUP_MEMBER_CACHE+uidstr+"_"+gidstr, load)
	if err != nil {
		return nil, err
	}
	groupmember, ok := result.(*dao.GroupMember)
	if !ok {
		err = errors.New("convert error")
	}
	return
}
func SetGroupMember(gid int, uid int, updateCol map[string]interface{}) (err error) {
	// 1.写入db
	// 2.删除缓存，包括 gid member list的缓存

	uidstr := strconv.Itoa(uid)
	gidstr := strconv.Itoa(gid)

	err = data_persistence.SaveGroupMember(uid, gid, updateCol)
	if err != nil {
		return err
	}

	err = data_cache.UnsetGroupMemberWithCache(uidstr, gidstr) //删除
	if err != nil {
		return err
	}

	err = data_cache.UnsetGroupMemberListWithCache(gidstr) //删除
	if err != nil {
		return err
	}

	err = data_cache.UnsetGroupListWithCache(uidstr) //删除
	if err != nil {
		return err
	}

	return err
}
func GroupMemberInvite(gid int, uids []int) (err error) {
	// 查询 group 信息
	// 查询 users 信息
	// 创建groupmember
	// todo 创建对应的coversation信息记录offset

	groupInfo, err := GetGroupInfo(gid)
	if err != nil {
		return
	}

	userList, err := GetUserList(uids)
	if err != nil {
		return
	}

	groupMemberList := make([]dao.GroupMember, 0, 16)
	for _, u := range userList {
		user := dao.GroupMember{
			GroupId:     groupInfo.Id,
			UserId:      u.Id,
			AliasName:   u.UserName,
			Remark:      groupInfo.GroupName,
			CreatedTime: time.Now(),
			IsDelete:    0,
		}
		groupMemberList = append(groupMemberList, user)
	}

	err = data_persistence.CreateGroupMemberInBatch(groupMemberList)
	if err != nil {
		return err
	}

	// todo 考虑涉及的数据，删除缓存

	err = data_cache.UnsetGroupMemberListWithCache(strconv.Itoa(gid))
	if err != nil {
		return err
	}

	return
}

/* EXT */
func GetGroupListByUser(uid int) (groupMemberList []dao.GroupMember, err error) {
	// 查询

	// 1.查询到数据
	// 2.查询错误
	// 3.查询不到数据
	//	查询数据库，放入缓存

	uidstr := strconv.Itoa(uid)

	groupMemberList, err = data_cache.GetGroupListWithCache(uidstr)
	if err == nil {
		return groupMemberList, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		gl, e := data_persistence.LoadGroupListByUId(uid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if len(gl) > 0 {
			e = data_cache.SetGroupListWithCache(uidstr, gl)
			if e != nil {
				return nil, e
			}
			return gl, nil
		}
		// 没有就设置空
		if len(gl) == 0 {
			e = data_cache.SetGroupListEmptyWithCache(uidstr)
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return gl, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_GROUP_LIST_CACHE+uidstr, load) // 给key加个前缀以防相同id
	if err != nil {
		return nil, err
	}
	groupMemberList, ok := result.([]dao.GroupMember)
	if !ok {
		err = errors.New("convert error")
	}
	return
}
func GetGroupList(gids []int) (groupList []dao.Group, err error) {

	groupList = make([]dao.Group, 0, 16)
	for _, gid := range gids {
		g, e := GetGroupInfo(gid)
		if e != nil {
			return // todo 部分遇到错误是否继续 ？continue
		}
		groupList = append(groupList, *g)
	}
	return
}
func GetGroupMemberList(gid int) (groupMemberList []dao.GroupMember, err error) {
	// 查询

	// 1.查询到数据
	// 2.查询错误
	// 3.查询不到数据
	//	查询数据库，放入缓存

	gidstr := strconv.Itoa(gid)

	groupMemberList, err = data_cache.GetGroupMemberListWithCache(gidstr)
	if err == nil {
		return groupMemberList, nil
	}
	if err != data_cache.ErrNoData { // 正常错误
		return nil, err
	}
	// 逻辑走到这里说明找不到数据并且redis没错误。查询数据库
	load := func() (interface{}, error) {
		gml, e := data_persistence.LoadGroupMemberList(gid)
		// 系统级错误，返回错误
		if e != nil {
			return nil, e
		}
		// 找到对应的数据item 就放入缓存
		if len(gml) > 0 {
			e = data_cache.SetGroupMemberListWithCache(gidstr, gml)
			if e != nil {
				return nil, e
			}
			return gml, nil
		}
		// 没有就设置空
		if len(gml) == 0 {
			e = data_cache.SetGroupMemberListEmptyWithCache(gidstr)
			if e != nil {
				return nil, e
			}
			// todo data_manager 没有数据是否需要 返回一个没有数据的err
			return gml, nil
		}
		return nil, e
	}
	result, err, _ := _singleFlight.Do(data_cache.PREFIX_GROUP_MEMBER_LIST_CACHE+gidstr, load)
	if err != nil {
		return nil, err
	}
	groupMemberList, ok := result.([]dao.GroupMember)
	if !ok {
		err = errors.New("convert error")
	}

	return
}

// 有些事务，分开的crud不能满足，提供一个专门的接口处理
func AddGroup(uid int, aliasName string, groupName string) (group *dao.Group, groupMember *dao.GroupMember, err error) {
	// 创建group

	uidstr := strconv.Itoa(uid)

	group, groupMember, err = data_persistence.CreateGroup(uid, aliasName, groupName)
	if err != nil {
		return nil, nil, err
	}

	err = data_cache.UnsetGroupListWithCache(uidstr) //删除
	if err != nil {
		return nil, nil, err
	}

	// todo 考虑涉及的缓存删除
	return
}
func DisbandGroup(gid int) (err error) {
	// 1.写入db
	// 2.删除缓存

	gidstr := strconv.Itoa(gid)

	err = data_persistence.DistoryGroup(gid)
	if err != nil {
		return err
	}

	err = data_cache.UnsetGroupWithCache(gidstr) //删除 group
	if err != nil {
		return err
	}

	gml, err := data_cache.GetGroupMemberListWithCache(gidstr)
	if err != nil {
		return err
	}
	for _, gm := range gml {
		data_cache.UnsetGroupListWithCache(strconv.Itoa(gm.UserId))           // 删除 这个group 下的成员的 groupList
		data_cache.UnsetGroupMemberWithCache(strconv.Itoa(gm.UserId), gidstr) // 删除 这个group 下的成员的 groupMember 信息
	}

	err = data_cache.UnsetGroupMemberListWithCache(gidstr) // 删除这个group下的成员信息
	if err != nil {
		return err
	}

	return
}
