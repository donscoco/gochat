package controller

import (
	"encoding/json"
	"fmt"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/conn_manager"
	"github.com/donscoco/gochat/internal/module/ws_sender"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{ // 用于给每次http连接设置升级为长链接websocket
	CheckOrigin: func(r *http.Request) bool { //解决跨域问题,不检查，都返回true
		return true
	},
}

type WSController struct{}

func WSRegister(group *gin.RouterGroup) {
	wsController := &WSController{}
	group.GET("/ws", wsController.IM)
}
func (wc *WSController) IM(c *gin.Context) {

	// 获取信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// 扩展为长连接websocket
	//返回头
	responseHeader := http.Header{}
	//responseHeader.Add("Sec-WebSocket-Protocol", "protoo")
	//升级为长连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, responseHeader)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	// 创建channel，用来和conn通信
	recvChan := make(chan []byte) // 接收来自browser的数据
	sendChan := make(chan []byte) // 接收发给browser的数据

	ws_sender.SetChan(strconv.Itoa(sessionInfo.ID), sendChan)

	// fordebug
	//go func() {
	//	for {
	//		time.Sleep(20 * time.Second)
	//		sendChan <- []byte(`{"cmd":0,"data":{"userId":1}}`)
	//	}
	//}()

	// 注册到 连接注册中心
	conn_manager.KeepAlive(sessionInfo.ID)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				// 关闭下游
				close(recvChan)
				return
			}
			recvChan <- message
		}
	}()

	ticker := time.NewTicker(10 * time.Second)
	lastHeartbeat := time.Now()
	for {

		select {
		//case  // 服务退出的通知，用于安全退出 暂时先不管
		case msg, ok := <-recvChan:
			if !ok {
				// todo 连接已关闭（客户端主动关闭）
				return
			}

			//log.Println("read:", string(msg)) // for debug

			// conn 存活设置
			conn_manager.KeepAlive(sessionInfo.ID)

			err = conn.WriteMessage(websocket.TextMessage, msg) // 默认是 text类型的消息
			if err != nil {
				// todo 这里发送的错误可能是超时，考虑是直接断开连接还是重试？
				return
			}

			lastHeartbeat = time.Now()
		case msg, ok := <-sendChan: // 放一个用于发送的channel
			if !ok {
				// todo 连接已关闭（服务端主动关闭）
				return
			}

			log.Println("send:", string(msg)) // for debug

			err = conn.WriteMessage(websocket.TextMessage, msg) // 默认是 text类型的消息
			if err != nil {
				// todo 这里发送的错误可能是超时，考虑是直接断开连接还是重试？
				return
			}

		case <-ticker.C: // 超时检查
			if lastHeartbeat.Add(10*time.Second).Unix() < time.Now().Unix() {
				return
			}
			log.Println("check heartbeat")
		}
	}
}
