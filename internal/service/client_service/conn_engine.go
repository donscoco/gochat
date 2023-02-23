package client_service

import (
	"github.com/donscoco/gochat/internal/module/coordinator/conn_coordinator"
	"github.com/donscoco/gochat/internal/module/rpc_client"
	"github.com/donscoco/gochat/pkg/iron_config"
)

var DefaultConnEngineCli *rpc_client.RpcClient

func InitConnEngineRpcClient() (err error) {

	addrs := make([]string, 0, 16)
	iron_config.Conf.GetByScan("/coordinator/zookeeper/addrs", &addrs)

	rootPath := iron_config.Conf.GetString("/rpc_server/conn_engine/root_path")

	err = conn_coordinator.InitZKCoordinator(addrs, rootPath)
	if err != nil {
		return err
	}

	cli, err := rpc_client.CreateRpcClient(
		3, 2000,
		conn_coordinator.Watcher,
		conn_coordinator.UpdateAddrList,
	)
	if err != nil {
		return err
	}

	err = cli.Start()
	if err != nil {
		return err
	}

	DefaultConnEngineCli = cli

	return nil
}

// todo 退出的时候，注销掉coordinator，虽然zk临时节点断开会自己关闭
func DestroyConnEngineRpcClient() (err error) {
	return nil
}
