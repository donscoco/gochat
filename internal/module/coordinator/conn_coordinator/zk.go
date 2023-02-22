package conn_coordinator

import (
	"encoding/json"
	"fmt"
	"github.com/donscoco/gochat/internal/module/rpc_client"
	"github.com/donscoco/gochat/internal/module/rpc_server"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

// 数据引擎是用的 rpc
var zkc *ZookeeperCoordinator

type ZookeeperCoordinator struct { // 一个rpc对应一个
	RootPath string   // 服务所属
	Value    State    // 用来保存每个节点的值，对应就是每个rpc服务的状态信息
	Stat     *zk.Stat // 用来保存zk节点的 stat
	Conn     *zk.Conn
}

// 服务节点信息的数据结构
type State struct {
	Name string
	Addr string

	Data string

	StartTime int64 // 服务开始时间
	Heartbeat int64 // 最后的心跳时间
}

func InitZKCoordinator(addrs []string, rootPath string) (err error) {
	if zkc != nil {
		return nil
	}
	zkc = new(ZookeeperCoordinator)
	zkc.RootPath = rootPath
	zkc.Conn, _, err = zk.Connect(addrs, 10*time.Second)
	if err != nil {
		return err
	}
	return

}

// 注册到注册中心
func Reigster(server *rpc_server.RpcServer) (err error) {

	// 创建服务节点，加上判断服务节点是否存在
	isExists, _, err := zkc.Conn.Exists(zkc.RootPath)
	if err != nil {
		return err
	}
	if !isExists { // 不存在就创建（永久节点）
		_, err = zkc.Conn.Create(zkc.RootPath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	// 创建server对应的临时节点
	info := State{
		Name:      "de_" + server.Addr, // todo 改成配置文件
		Addr:      server.Addr,
		Data:      "",
		StartTime: time.Now().Unix(),
		Heartbeat: 0,
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	str, err := zkc.Conn.Create(zkc.RootPath+"/"+info.Addr, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Println("test")
		return err
	}
	_, stat, err := zkc.Conn.Exists(zkc.RootPath + "/" + info.Addr) // 拿到stat
	if err != nil {
		return err
	}
	zkc.Value = info
	zkc.Stat = stat

	log.Println(str)

	return nil
}

// 发送心跳
func Heartbeat(server *rpc_server.RpcServer) (err error) {
	ticker := time.NewTicker(10 * time.Second) // todo 文件配置
	for {
		select {
		case <-server.Ctx.Done():
			// 收到退出通知，关闭下游
			return
		case <-ticker.C:

			zkc.Value.Heartbeat = time.Now().Unix()
			data, _ := json.Marshal(zkc.Value) // 前面marshal没有错，这里marshal肯定也不会错直接不用管err

			stat, err := zkc.Conn.Set(zkc.RootPath+"/"+server.Addr, data, zkc.Stat.Version)
			if err != nil {
				// todo 重试次数超时报警等机制
				// 连接不上coordinator，尝试重新注册
				err = server.Register(server)
				if err != nil {
					fmt.Println(err)
					//return
				}
			}
			zkc.Stat = stat
		}
	}
	return nil
}

// watch 节点更新
func Watcher(client *rpc_client.RpcClient) (err error) {

	// watch 节点数据变化 节点删除 节点的子节点变化

	var stat *zk.Stat

	var ExistEvent <-chan zk.Event
	var DataEvent <-chan zk.Event
	var ChildEvent <-chan zk.Event
	var exist bool
	var data []byte
	var child []string

	exist, stat, ExistEvent, err = zkc.Conn.ExistsW(zkc.RootPath)
	if err != nil {
		return err
	}

	data, stat, DataEvent, err = zkc.Conn.GetW(zkc.RootPath)
	if err != nil {
		return err
	}

	child, stat, ChildEvent, err = zkc.Conn.ChildrenW(zkc.RootPath)
	if err != nil {
		return err
	}

	zkc.Stat = stat

	// 开启监听前，先同步下

	for {
		select {
		case e := <-ExistEvent:
			// 检查
			if e.Type != zk.EventNodeCreated || e.Type != zk.EventNodeDeleted {
				// todo
			}

			// todo
			log.Println("是否存在节点:", exist)

			// watch
			exist, stat, ExistEvent, err = zkc.Conn.ExistsW(zkc.RootPath)
			if err != nil {
				return err
			}
			zkc.Stat = stat
		case e := <-DataEvent:
			if e.Type != zk.EventNodeDataChanged {
				// todo
			}

			// todo
			log.Println("节点数据:", string(data))

			data, stat, DataEvent, err = zkc.Conn.GetW(zkc.RootPath)
			if err != nil {
				return nil
			}
			zkc.Stat = stat
		case e := <-ChildEvent:
			// 检查
			if e.Type != zk.EventNodeChildrenChanged {
				// todo
			}

			// todo
			//获取最新的服务uri
			// 孩子节点出现变化，获取所有孩子节点，创建对应信息，
			err := client.UpdateClient(client)
			if err != nil {
				return err
			}
			log.Println("孩子节点:", child)

			// watch
			child, stat, ChildEvent, err = zkc.Conn.ChildrenW(zkc.RootPath)
			if err != nil {
				return err
			}
			zkc.Stat = stat

		case <-client.Context.Done():
			// 关闭下游
			return
		}
	}

	return nil
}

// 更新client 的维护的服务地址
func UpdateAddrList(client *rpc_client.RpcClient) error {
	newAddrs := make([]string, 0, 16)
	child, stat, err := zkc.Conn.Children(zkc.RootPath)
	//cNodes, err := c.serverNode.GetChildrenNodes()
	if err != nil {
		return err
	}
	zkc.Stat = stat
	for _, cn := range child {
		data, _, err := zkc.Conn.Get(zkc.RootPath + "/" + cn)
		if err != nil {
			return err
		}
		serverInfo := &State{}
		err = json.Unmarshal(data, serverInfo)
		if err != nil {
			return err
		}
		newAddrs = append(newAddrs, serverInfo.Addr)
	}

	//更新agent
	newAgent := make(map[string]*rpc_client.ClientAgent)
	newBalance := make([]*rpc_client.ClientAgent, 0, 16)
	for _, addr := range newAddrs {
		agentp, exist := client.ClientAgent[addr]
		if exist {
			newAgent[addr] = agentp
			newBalance = append(newBalance, agentp)
		} else { // 没有就创建新的agent
			agentp, err = rpc_client.CreateClientAgent(addr, 2000, 10, 2) // todo 文件配置
			if err != nil {
				return err
			}
			newAgent[addr] = agentp
			newBalance = append(newBalance, agentp)
		}
	}
	oldAgent := client.ClientAgent
	client.ClientAgent = newAgent
	client.Balance = newBalance

	//关闭旧agent
	for addr, agent := range oldAgent {
		_, exist := newAgent[addr]
		if exist {
			continue
		}
		agent.Close()
	}

	return nil
}
