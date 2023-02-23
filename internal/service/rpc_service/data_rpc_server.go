package rpc_service

import (
	"github.com/donscoco/gochat/internal/handler/data_handler"
	"github.com/donscoco/gochat/internal/module/coordinator/data_coordinator"
	"github.com/donscoco/gochat/internal/module/rpc_server"
	"github.com/donscoco/gochat/pkg/iron_config"
)

var dataEngineRpcServer *rpc_server.RpcServer

func DataEngineRpcServerStart() (err error) {

	addrs := make([]string, 0, 16)
	iron_config.Conf.GetByScan("/coordinator/zookeeper/addrs", &addrs)

	rootPath := iron_config.Conf.GetString("/rpc_server/data_engine/root_path")
	rpcServerAddr := iron_config.Conf.GetString("/rpc_server/data_engine/addr")

	err = data_coordinator.InitZKCoordinator(addrs, rootPath)
	if err != nil {
		return err
	}

	rpcserver, err := rpc_server.CreateRpcServer(
		//":9090",
		rpcServerAddr,
		&data_handler.DataEngine{},
		data_coordinator.Reigster,
		data_coordinator.Heartbeat,
	)
	if err != nil {
		return err
	}

	err = rpcserver.Start()
	if err != nil {
		return err
	}
	dataEngineRpcServer = rpcserver

	return nil

}
func DataEngineRpcServerStop() error {
	return dataEngineRpcServer.Stop()
}
