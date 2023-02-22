package main

import (
	"flag"
	"github.com/donscoco/gochat/config"
	"github.com/donscoco/gochat/internal/service"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/donscoco/gochat/internal/service/msg_service/msg_producer"
	"github.com/donscoco/gochat/internal/service/rpc_service"
	"github.com/donscoco/gochat/pkg/iron_config"
	"github.com/donscoco/gochat/pkg/iron_core"
	"github.com/donscoco/gochat/pkg/iron_log"
)

var env string

var GoCore *iron_core.Core
var Config *iron_config.Config
var Logger iron_core.Logger

func init() {
	flag.StringVar(&env, "env", "dev", "")
}

func main() {

	flag.Parse()

	configPath := config.Path("./" + env + "/data_engine/config.json")

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
		//core.InitDBManager,
		service.InitRedisService,
		service.InitMysqlService,
		service.InitGORMService,
		service.InitMgoService,
		rpc_service.DataEngineRpcServerStart,
		client_service.InitConnEngineRpcClient,
		msg_producer.MsgProducerStart,
	)

	GoCore.OnStop( //先关上游
		msg_producer.MsgProducerStop,
		//rpc_service.DataEngineRpcServerStop, // 酌情考量，因为rpc client 是维护着长链接的。如果加上，安全退出就会一直等待上游rpc client先关闭
		service.StopGORMService,
		service.StopMysqlService,
		service.StopRedisService,
	)

	GoCore.Boot()

	iron_log.Close()

}
