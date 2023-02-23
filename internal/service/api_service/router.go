package api_service

import (
	"github.com/donscoco/gochat/internal/controller"
	"github.com/donscoco/gochat/internal/middleware/http_middleware"
	"github.com/donscoco/gochat/pkg/iron_config"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"

	conf "github.com/donscoco/gochat/config"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	config := iron_config.Conf
	router := gin.Default()

	// debug
	router.Use(middlewares...)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	store, err := sessions.NewRedisStore(10, "tcp",
		config.GetString("/session/redis_addr"),
		config.GetString("/session/redis_pwd"),
		[]byte("secret"),
	)
	if err != nil {
		//log.Errorf("sessions.NewRedisStore err:%v", err)
		log.Fatalf("sessions.NewRedisStore err:%v", err)
	}

	// 登陆注册接口
	LoginRouter := router.Group("/api")
	LoginRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.TranslationMiddleware(),
	)
	{
		controller.LoginRegister(LoginRouter)
	}

	// 用户接口
	userRouter := router.Group("/api/user")
	userRouter.Use(
		sessions.Sessions("mysession", store), // 设置cookie 名
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.UserRegister(userRouter)
	}

	// 好友接口
	contactsRouter := router.Group("/api/friend")
	contactsRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.ContactsRegister(contactsRouter)
	}

	// 群接口
	groupRouter := router.Group("/api/group")
	groupRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.GroupRegister(groupRouter)
	}

	// 媒体接口
	mediaRouter := router.Group("/api")
	mediaRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.MediaRegister(mediaRouter)
	}

	// 消息接口
	msgRouter := router.Group("/api/message")
	msgRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.MsgRegister(msgRouter)
	}

	// ws接口
	wsRouter := router.Group("/im")
	wsRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.WSRegister(wsRouter)
	}

	//  webRTC
	iceRouter := router.Group("/api/webrtc")
	iceRouter.Use(
		sessions.Sessions("mysession", store),
		http_middleware.RecoveryMiddleware(),
		http_middleware.RequestLog(),
		http_middleware.SessionAuthMiddleware(), // 检查是否已经有登陆的session
		http_middleware.TranslationMiddleware())
	{
		controller.WebRTCRegister(iceRouter)
	}

	// 前端静态资源
	//router.Static("/dist", "./dist")
	router.Static("/dist", conf.Path("../dist"))

	return router
}
