package main

import (
	"github.com/donscoco/gochat/internal/module/ws_server"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ws_server.DefaultWebsocketServer = ws_server.NewWebsocketServer(
		":8878",
		connHandle,
	)

	ws_server.DefaultWebsocketServer.Start()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ws_server.DefaultWebsocketServer.Stop()
}

func connHandle(socket *websocket.Conn, request *http.Request) {

}
