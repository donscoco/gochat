package data_cache

import (
	"encoding/json"
	"errors"
	"github.com/donscoco/gochat/internal/base/redis"
	"github.com/donscoco/gochat/internal/dao"
	"log"
)

// set
//		获取redis conn
//		获取key
// get
//
// scan

func GetUserWithCache(uid string) (user *dao.User, err error) {
	// key  user_id
	// 获取连接
	// 有就返回，没有就sigleflight

	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		//todo
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_USER_CACHE + uid)

	if cmd != nil && cmd.Err() == nil {
		user = new(dao.User)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的 user
			// todo
			return user, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), user) // todo 这里用不用初始化user来着？
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [user]")

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

func SetUserWithCache(uid string, user *dao.User) (err error) {

	val, err := json.Marshal(user)
	if err != nil {
		// todo
		return err
	}

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_USER_CACHE+uid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func SetUserEmptyWithCache(uid string) (err error) {
	// key  user_id
	// 获取连接

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_USER_CACHE+uid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func UnsetUserWithCache(uid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_USER_CACHE + uid)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	log.Println("domark 缓存删除 [user]")

	return nil

}
