package main

import (
	"flag"
	"github.com/donscoco/gochat/config"
	"github.com/donscoco/gochat/internal/service"
	"github.com/donscoco/gochat/internal/service/api_service"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/donscoco/gochat/internal/service/rpc_service"
	"github.com/donscoco/gochat/pkg/iron_config"
	"github.com/donscoco/gochat/pkg/iron_core"
	"github.com/donscoco/gochat/pkg/iron_log"
	"os"
)

const GOCHAT_ENV = "GOCHAT_ENV"

var env string

var GoCore *iron_core.Core
var Config *iron_config.Config
var Logger iron_core.Logger

func init() {
	flag.StringVar(&env, "env", "dev", "")
}

// swag init -g cmd/conn_engine/main.go -o internal/swagger/conn_engine
func main() {

	flag.Parse()

	gochatenv := os.Getenv(GOCHAT_ENV)
	if len(gochatenv) > 0 {
		env = gochatenv
	}

	configPath := config.Path("./" + env + "/conn_engine/config.json")

	// 初始化config
	iron_config.Conf = iron_config.NewConfiguration(configPath)

	// 初始化log
	//iron_log.InitLoggerByEnv()
	iron_log.InitLoggerByParam(
		iron_config.Conf.GetString("/log/log_path"),
		iron_config.Conf.GetString("/log/log_level"),
		iron_config.Conf.GetString("/log/log_mode"),
	)

	// 初始化core
	Logger = iron_log.NewLogger("CORE")
	GoCore = iron_core.NewCore()
	GoCore.SetLogger(Logger)

	GoCore.OnStart(
		service.InitOSSService,
		service.InitRedisService,
		service.InitMysqlService,
		service.InitGORMService,

		service.InitContainer, //用来存放 发给websocket的channel

		client_service.InitDataEngineRpcClient,
		rpc_service.ConnEngineRpcServerStart,
		api_service.HttpServerRun,
	)

	GoCore.OnStop( //先关上游
		api_service.HttpServerStop,
		rpc_service.ConnEngineRpcServerStop,

		service.StopGORMService,
		service.StopMysqlService,
		service.StopRedisService,
	)

	GoCore.Boot()

	iron_log.Close()

}
