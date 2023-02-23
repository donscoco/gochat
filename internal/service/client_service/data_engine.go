package client_service

import (
	"github.com/donscoco/gochat/internal/module/coordinator/data_coordinator"
	"github.com/donscoco/gochat/internal/module/rpc_client"
	"github.com/donscoco/gochat/pkg/iron_config"
)

var DefaultDataEngineCli *rpc_client.RpcClient

func InitDataEngineRpcClient() (err error) {

	addrs := make([]string, 0, 16)
	iron_config.Conf.GetByScan("/coordinator/zookeeper/addrs", &addrs)

	rootPath := iron_config.Conf.GetString("/rpc_server/data_engine/root_path")

	err = data_coordinator.InitZKCoordinator(addrs, rootPath)
	if err != nil {
		return err
	}

	cli, err := rpc_client.CreateRpcClient(
		3, 2000,
		data_coordinator.Watcher,
		data_coordinator.UpdateAddrList,
	)
	if err != nil {
		return err
	}

	err = cli.Start()
	if err != nil {
		return err
	}

	DefaultDataEngineCli = cli

	return nil
}
