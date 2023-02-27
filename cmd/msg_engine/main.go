package main

import (
	"flag"
	"github.com/donscoco/gochat/config"
	"github.com/donscoco/gochat/internal/service"
	"github.com/donscoco/gochat/internal/service/client_service"
	"github.com/donscoco/gochat/internal/service/msg_service/msg_consumer"
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
func main() {

	flag.Parse()

	gochatenv := os.Getenv(GOCHAT_ENV)
	if len(gochatenv) > 0 {
		env = gochatenv
	}

	configPath := config.Path("./" + env + "/msg_engine/config.json")

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
		service.InitRedisService, //使用websocket的注册中心，寻找 文件描述符 所在机器
		service.InitMysqlService,
		service.InitGORMService,
		service.InitMgoService,
		client_service.InitDataEngineRpcClient,
		client_service.InitConnEngineRpcClient,
		msg_consumer.MsgConsumerStart,
	)
	GoCore.OnStop( //先关上游
		msg_consumer.MsgConsumerStop,
		service.StopGORMService,
		service.StopMysqlService,
		service.StopRedisService,
	)
	GoCore.Boot()

	iron_log.Close()
}
