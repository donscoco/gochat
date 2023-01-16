package main

import (
	"fmt"
	_ "github.com/gorilla/websocket" // 测试
	"net"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", connectHandler)

	server := &http.Server{
		Addr:         ":30010",
		WriteTimeout: time.Second * 10,
		Handler:      mux,
	}

	fmt.Print("test")

	server.ListenAndServe()
}

func connectHandler(resp http.ResponseWriter, req *http.Request) {

	// 检查ip
	resp.Write([]byte("<<<<<<<<<<<<<<<< check ip >>>>>>>>>>>>>>>>\n"))
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		resp.Write([]byte(err.Error() + "\n"))
	}

	for _, addr := range addrs {
		resp.Write([]byte(addr.String() + "\n"))
		//log.Println(addrs)
	}
}
