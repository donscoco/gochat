package ws_server

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var DefaultWebsocketServer *WebsocketServer

type WebsocketServer struct {
	// 放一个处理函数
	Addr            string
	handleWebSocket func(socket *websocket.Conn, request *http.Request) // 用于服务启动后处理每个连接的方法
	upgrader        websocket.Upgrader

	httpServer http.Server

	ctx    context.Context
	cancel context.CancelFunc

	sync.WaitGroup
	sync.Mutex
}

func NewWebsocketServer(addr string, wsHandler func(socket *websocket.Conn, request *http.Request)) (ws *WebsocketServer) { // todo 注册一个处理函数，注册一个关闭函数进来
	ws = new(WebsocketServer)
	ws.Addr = addr
	ws.handleWebSocket = wsHandler
	ws.httpServer = http.Server{
		Addr: ws.Addr,
		//Handler: mux,
		//TLSConfig:
	}

	ws.upgrader = websocket.Upgrader{ // 用于给每次http连接设置升级为长链接websocket
		CheckOrigin: func(r *http.Request) bool { //解决跨域问题,不检查，都返回true
			return true
		},
	}
	return
}

func (ws *WebsocketServer) Start() {

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.handleWebSocketRequest)
	mux.HandleFunc("/echo", ws.echo) // for debug
	mux.HandleFunc("/im", ws.echo)   // for debug

	ws.httpServer.Handler = mux

	go func() {
		log.Println("listen at:", ws.Addr)
		err := ws.httpServer.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}
func (ws *WebsocketServer) Stop() error {
	ws.cancel()
	// 关闭各种资源
	//ws.httpServer.Close()
	ws.Wait()
	log.Println("安全退出")
	return nil
}

// 处理函数
func (ws *WebsocketServer) handleWebSocketRequest(writer http.ResponseWriter, request *http.Request) {
	//返回头
	responseHeader := http.Header{}
	//responseHeader.Add("Sec-WebSocket-Protocol", "protoo")
	//升级为长连接
	conn, err := ws.upgrader.Upgrade(writer, request, responseHeader)
	//输出错误日志
	if err != nil {
		// todo
		//util.Panicf("%v", err)
		log.Println("%v", err)
	}
	go ws.serve(conn, request)
}

func (ws *WebsocketServer) serve(conn *websocket.Conn, request *http.Request) {
	defer conn.Close()
	ws.Add(1)                         // todo 考虑一下是否真的要这样做安全退出，看别人都是直接不管conn的。
	ws.handleWebSocket(conn, request) // 使用创建的时候注册的处理函数处理
	ws.Done()                         //
}

func (ws *WebsocketServer) echo(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		log.Printf("%s recv: %s", ws.Addr, message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
