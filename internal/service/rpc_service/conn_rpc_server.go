package rpc_service

import (
	"github.com/donscoco/gochat/internal/handler/conn_handler"
	"github.com/donscoco/gochat/internal/module/coordinator/conn_coordinator"
	"github.com/donscoco/gochat/internal/module/rpc_server"
	"github.com/donscoco/gochat/pkg/iron_config"
)

var connEngineRpcServer *rpc_server.RpcServer

func ConnEngineRpcServerStart() (err error) {

	addrs := make([]string, 0, 16)
	iron_config.Conf.GetByScan("/coordinator/zookeeper/addrs", &addrs)

	rootPath := iron_config.Conf.GetString("/rpc_server/conn_engine/root_path")
	rpcServerAddr := iron_config.Conf.GetString("/rpc_server/conn_engine/addr")

	err = conn_coordinator.InitZKCoordinator(addrs, rootPath)
	if err != nil {
		return err
	}

	rpcserver, err := rpc_server.CreateRpcServer(
		//":9080",
		rpcServerAddr,
		&conn_handler.ConnEngine{},
		conn_coordinator.Reigster,
		conn_coordinator.Heartbeat,
	)
	if err != nil {
		return err
	}

	err = rpcserver.Start()
	if err != nil {
		return err
	}
	connEngineRpcServer = rpcserver

	return nil

}
func ConnEngineRpcServerStop() error {
	return connEngineRpcServer.Stop()
}
