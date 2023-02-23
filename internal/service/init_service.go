package service

import (
	"errors"
	"github.com/donscoco/gochat/internal/base/mongodb"
	db "github.com/donscoco/gochat/internal/base/mysql"
	"github.com/donscoco/gochat/internal/base/oss"
	"github.com/donscoco/gochat/internal/base/redis"
	"github.com/donscoco/gochat/internal/module/ws_sender"
	"github.com/donscoco/gochat/pkg/gorm"
	"github.com/donscoco/gochat/pkg/iron_config"
)

// / 初始化 GORM
func InitGORMService() error {
	//gorm.DefaultDB = core.GoCore.GetConf().GetString("/core/dbManager/mysql/0/proxyname")
	//// gorm
	//mysqlProxy, err := core.GoCore.DB.GetMySQL(gorm.DefaultDB)
	//if err != nil {
	//	return err
	//}
	//gorm.InitDBPool(
	//	gorm.DefaultDB,
	//	mysqlProxy.Session,
	//)
	//
	//g, err := gorm.GetGormPool(gorm.DefaultDB)
	//fmt.Println(g.Statement)
	//return nil

	gorm.DefaultDB = iron_config.Conf.GetString("/server/mysql/0/proxy_name")
	// gorm
	mysqlProxy, ok := db.MySQLManager[gorm.DefaultDB]
	if !ok {
		return errors.New("init gorm fail, empty mysql")
	}
	err := gorm.InitDBPool(
		gorm.DefaultDB,
		mysqlProxy.Session,
	)
	if err != nil {
		return err
	}
	return nil
}

func StopGORMService() error {
	return gorm.CloseDB()
}

// / 初始化 OSS
func InitOSSService() error {

	oss.DefaultEndpoint = iron_config.Conf.GetString("/oss/endpoint")
	oss.DefaultAccessKeyId = iron_config.Conf.GetString("/oss/access_key_id")
	oss.DefaultAccessKeySecret = iron_config.Conf.GetString("/oss/access_key_secret")
	oss.DefaultBucket = iron_config.Conf.GetString("/oss/bucket")
	oss.DefaultURLPrefix = iron_config.Conf.GetString("/oss/url_prefix")

	return oss.InitOSS(
		oss.DefaultEndpoint,
		oss.DefaultAccessKeyId,
		oss.DefaultAccessKeySecret,
	)
}

// 初始化存放 websocket 的 channel
func InitContainer() error {
	return ws_sender.InitConnChanMap() // map中拿到的channel 能给ws发消息
}

func InitMgoService() error {

	addrs := make([]string, 0, 16)
	iron_config.Conf.GetByScan("/server/mongodb/0/addrs", &addrs)

	var conntimeout int
	iron_config.Conf.GetByScan("/server/mongodb/0/connect_timeout", &conntimeout)

	return mongodb.InitMgoClient(addrs, conntimeout)
}

// 初始化 mysql
func InitMysqlService() error {

	proxys := make([]db.MySQLProxy, 0, 16)
	iron_config.Conf.GetByScan("/server/mysql", &proxys)

	return db.InitMySQLProxys(proxys)
}

func StopMysqlService() error {
	for _, p := range db.MySQLManager {
		err := p.Session.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func InitRedisService() error {

	proxys := make([]redis.RedisConfig, 0, 16)
	iron_config.Conf.GetByScan("/server/redis", &proxys)

	return redis.InitSingleRedisClient(proxys)
}

func StopRedisService() error {
	for _, p := range redis.RedisSingleClients {
		err := p.Close()
		if err != nil {
			return err
		}
	}

	for _, p := range redis.RedisClusterClients {
		err := p.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
