package data_cache

import (
	"encoding/json"
	"errors"
	"github.com/donscoco/gochat/internal/base/redis"
	"github.com/donscoco/gochat/internal/dao"
	"log"
)

func GetContactsWithCache(uid string, fid string) (contact *dao.Contacts, err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		// 这里是否要打印日志，还是统一往上抛，到服务层去记录日志？
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_CONTACTS_CACHE + uid + "_" + fid)

	if cmd != nil && cmd.Err() == nil {
		contact = new(dao.Contacts)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的
			// todo
			return contact, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), contact)
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [contacts] ")

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

func SetContactsWithCache(uid string, fid string, contact *dao.Contacts) (err error) {

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

	cmd := p.Set(PREFIX_CONTACTS_CACHE+uid+"_"+fid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func SetContactsEmptyWithCache(uid string, fid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_CONTACTS_CACHE+uid+"_"+fid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}

func UnsetContactsWithCache(uid string, fid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_CONTACTS_CACHE + uid + "_" + fid)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	log.Println("domark 缓存删除 [contacts] ")

	return nil

}

// ext

func GetContactsListWithCache(uid string) (contacts []dao.Contacts, err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return nil, errors.New("empty redis")
	}

	cmd := p.Get(PREFIX_CONTACTS_LIST_CACHE + uid)

	if cmd != nil && cmd.Err() == nil {
		contacts = make([]dao.Contacts, 0, 16)
		if cmd.Val() == EMPTY { // 为了防止缓存穿透,不用访问db，直接返回 空的
			// todo
			return contacts, nil
		}
		err = json.Unmarshal([]byte(cmd.Val()), &contacts)
		if err != nil {
			return nil, err
		}

		log.Println("domark 缓存获取 [contactsList] ")

		return contacts, nil
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
func SetContactsListWithCache(uid string, contacts []dao.Contacts) (err error) {

	val, err := json.Marshal(contacts)
	if err != nil {
		// todo
		return err
	}

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_CONTACTS_LIST_CACHE+uid, string(val), TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}
func UnsetContactsListWithCache(uid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Del(PREFIX_CONTACTS_LIST_CACHE + uid)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}

	log.Println("domark 缓存删除 [contactsList] ")

	return nil

}
func SetContactsListEmptyWithCache(uid string) (err error) {

	p, ok := redis.RedisSingleClients["default"]
	if !ok { //todo log
		//todo
		return errors.New("empty redis")
	}

	cmd := p.Set(PREFIX_CONTACTS_LIST_CACHE+uid, EMPTY, TTL)

	if cmd != nil && cmd.Err() != nil {
		// todo
		return err
	}
	return nil

}
