package main

import (
	"flag"
	"github.com/donscoco/gochat/config"
	"github.com/donscoco/gochat/internal/service"
	"github.com/donscoco/gochat/internal/service/chatgpt_service"
	"github.com/donscoco/gochat/pkg/iron_config"
	"github.com/donscoco/gochat/pkg/iron_core"
	"github.com/donscoco/gochat/pkg/iron_log"
	"os"
)

// todo chatGPT 实现聊天机器人

const GOCHAT_ENV = "GOCHAT_ENV"

var env string

var GoCore *iron_core.Core
var Config *iron_config.Config
var Logger iron_core.Logger

func init() {
	flag.StringVar(&env, "env", "dev", "")
}

func main() {

	// 1.查询判断db，看是否创建聊天机器人
	// 2.http登陆，对于的session
	// 3.http keep alive websocket
	// 4.接收websocket的 msg 类型数据结构
	//
	//	4.1 text 以外的消息返回不支持
	//	4.2 发送到 client的 goroutine里面

	////////
	//先创建chatgpt的client

	flag.Parse()

	gochatenv := os.Getenv(GOCHAT_ENV)
	if len(gochatenv) > 0 {
		env = gochatenv
	}

	configPath := config.Path("./" + env + "/chat_bot/config.json")

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

	/// test  start

	//robot := make(map[string]interface{})
	service.InitMysqlService()
	service.InitGORMService()
	/// test end

	GoCore.OnStart(
		service.InitMysqlService,
		service.InitGORMService,
		chatgpt_service.Start,
	)
	GoCore.OnStop(
		chatgpt_service.Stop,
		service.StopGORMService,
		service.StopMysqlService,
	)
	GoCore.Boot()

	iron_log.Close()

}
