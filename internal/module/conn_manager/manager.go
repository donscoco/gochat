package conn_manager

import (
	"github.com/donscoco/gochat/internal/base/redis"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/pkg/iron_config"
	"github.com/donscoco/gochat/pkg/util"
	"log"
	"net"
	"strconv"
	"time"
)

// 提供 conn 和ip 注册到 fd 管理中心（这里简单做就用redis保存）
const (
	CONN_FD_PREFIX = "conn_fd_"
)

func KeepAlive(uid int) {
	// 获取连接
	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		//todo log
		bl.Error("[conn_manager] get redis fail")
		return
	}

	// 获取本机ip
	ip := util.GetLocalIp() // todo 考虑k8s，考虑获取机器名
	//hostname,err := os.Hostname()

	conf := iron_config.Conf
	addr := conf.GetString("/rpc_server/conn_engine/addr")
	_, port, _ := net.SplitHostPort(addr)

	cmd := p.Set(CONN_FD_PREFIX+strconv.Itoa(uid), ip+":"+port, 10*time.Second) // todo 文件配置，过期时间要考虑前端的存活探测频率
	if cmd != nil && cmd.Err() != nil {
		//return nil, cmd.Err()
		log.Println(cmd)
		// todo log
		return
	}

	return
}

func IsAlive(uid int) bool {
	// 获取连接
	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		// todo log
		return false
	}

	// todo 这里是否加个判断nil比较好
	cmd := p.Get(CONN_FD_PREFIX + strconv.Itoa(uid)) // todo 文件配置，过期时间要考虑前端的存活探测频率
	if cmd != nil && cmd.Err() != nil {
		//return nil, cmd.Err()
		log.Println(cmd)
		// todo log
		return false
	}

	return len(cmd.Val()) > 0
}

func Where(uid int) (addr string) {
	// 获取连接
	p, ok := redis.RedisSingleClients["default"]
	if !ok {
		// todo log
		return ""
	}

	cmd := p.Get(CONN_FD_PREFIX + strconv.Itoa(uid)) // todo 文件配置，过期时间要考虑前端的存活探测频率
	if cmd != nil && cmd.Err() != nil {
		//return nil, cmd.Err()
		log.Println(cmd)
		// todo log
		return ""
	}

	return cmd.Val()
}

func Dead(uid int) {

}
