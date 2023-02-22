package data_cache

import (
	"encoding/json"
	"errors"
	"github.com/donscoco/gochat/internal/base/redis"
	"github.com/donscoco/gochat/internal/dao"
	"log"
)

func GetGroupWithCache(gid string) (user *dao.Group, err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		// todo log
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_GROUP_CACHE + gid)

	if cmd != nil && cmd.Err() == nil {
		user = new(dao.Group)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的 user
			// todo
			return user, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), user)
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [group] ")

		return user, nil
	}

	// 找不到数据的错误
	if cmd != nil && cmd.Err() != nil && cmd.Err().Error() == "redis: nil" {
		return nil, ErrNoData

	}
	// 真的错误
	if cmd != nil && cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return nil, err

}

func SetGroupWithCache(gid string, group *dao.Group) (err error) {

	val, err := json.Marshal(group)
	if err != nil {
		// todo
		return err
	}

	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_CACHE+gid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func SetGroupEmptyWithCache(gid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_CACHE+gid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func UnsetGroupWithCache(gid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_GROUP_CACHE + gid)

	log.Println("domark 缓存删除 [group] ")

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	return nil

}

/////

func GetGroupMemberWithCache(uid string, gid string) (contact *dao.GroupMember, err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_GROUP_MEMBER_CACHE + uid + "_" + gid)

	if cmd != nil && cmd.Err() == nil {
		contact = new(dao.GroupMember)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的
			// todo
			return contact, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), contact)
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [groupMember] ")

		return contact, nil
	}

	// 找不到数据的错误
	if cmd != nil && cmd.Err() != nil && cmd.Err().Error() == "redis: nil" {
		return nil, ErrNoData

	}
	// 真的错误
	if cmd != nil && cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return nil, err

}

func SetGroupMemberWithCache(uid string, gid string, contact *dao.GroupMember) (err error) {

	val, err := json.Marshal(contact)
	if err != nil {
		// todo
		return err
	}

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_MEMBER_CACHE+uid+"_"+gid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func SetGroupMemberEmptyWithCache(uid string, gid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_MEMBER_CACHE+uid+"_"+gid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func UnsetGroupMemberWithCache(uid string, gid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_GROUP_MEMBER_CACHE + uid + "_" + gid)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	log.Println("domark 缓存删除 [groupMember] ")

	return nil

}

///// ext

func GetGroupMemberListWithCache(gid string) (groupMember []dao.GroupMember, err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_GROUP_MEMBER_LIST_CACHE + gid)

	if cmd != nil && cmd.Err() == nil {
		groupMember = make([]dao.GroupMember, 0, 16)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的
			// todo
			return groupMember, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), &groupMember)
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [groupMember list] ")

		return groupMember, nil
	}

	// 找不到数据的错误
	if cmd != nil && cmd.Err() != nil && cmd.Err().Error() == "redis: nil" {
		return nil, ErrNoData

	}
	// 真的错误
	if cmd != nil && cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return nil, err

}
func SetGroupMemberListWithCache(gid string, groupMember []dao.GroupMember) (err error) {

	val, err := json.Marshal(groupMember)
	if err != nil {
		// todo
		return err
	}

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_MEMBER_LIST_CACHE+gid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}
func UnsetGroupMemberListWithCache(gid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_GROUP_MEMBER_LIST_CACHE + gid)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	log.Println("domark 缓存删除 [groupMember list] ")

	return nil

}
func SetGroupMemberListEmptyWithCache(gid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_MEMBER_LIST_CACHE+gid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

// ext

func GetGroupListWithCache(uid string) (groupMember []dao.GroupMember, err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_GROUP_LIST_CACHE + uid)

	if cmd != nil && cmd.Err() == nil {
		groupMember = make([]dao.GroupMember, 0, 16)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的
			// todo
			return groupMember, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), &groupMember)
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [group list] ")

		return groupMember, nil
	}

	// 找不到数据的错误
	if cmd != nil && cmd.Err() != nil && cmd.Err().Error() == "redis: nil" {
		return nil, ErrNoData

	}
	// 真的错误
	if cmd != nil && cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return nil, err

}
func SetGroupListWithCache(uid string, groupMember []dao.GroupMember) (err error) {

	val, err := json.Marshal(groupMember)
	if err != nil {
		// todo
		return err
	}

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_LIST_CACHE+uid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}
func UnsetGroupListWithCache(uid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_GROUP_LIST_CACHE + uid)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	log.Println("domark 缓存获取 [group list] ")

	return nil

}
func SetGroupListEmptyWithCache(uid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_GROUP_LIST_CACHE+uid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}
