package rpc_server

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"sync"
	"time"
)

// 注册的信息
type State struct {
	Name string // 节点名称
	Addr string

	Data string

	StartTime int64 // 服务开始时间
	Heartbeat int64 // 最后的心跳时间
}
type RpcServer struct {
	Addr string
	//HeartbeatSec int    // 心跳间隔 让心跳函数去维护

	// 注册到zk的状态信息
	//State State

	// 服务监听
	listener net.Listener

	// 注册中心不一定要zk来实现，使用这样的方式，屏蔽具体实现
	Register  func(*RpcServer) error // 注册函数
	Heartbeat func(*RpcServer) error // 心跳函数

	Ctx    context.Context
	Cancel context.CancelFunc
	sync.WaitGroup
}

func CreateRpcServer(addr string, handler interface{},
	register func(server *RpcServer) error,
	heartbeat func(server *RpcServer) error) (s *RpcServer, err error) {

	s = new(RpcServer)
	s.Register = register
	s.Heartbeat = heartbeat

	err = rpc.Register(handler)
	if err != nil {
		return nil, err
	}

	// 如果没有指定ip则自己查找默认的子网ip 192.168 //fixme 兼容处理 k8s下的网络
	host, port, _ := net.SplitHostPort(addr) // :port
	if len(host) == 0 {
		host = GetIpv4_192_168()
	}
	if len(host) == 0 {
		host = GetIpv4_172()
	}
	if len(host) == 0 {
		host = GetIpv4_10()
	}
	s.Addr = host + ":" + port

	s.Ctx, s.Cancel = context.WithCancel(context.TODO())

	return
}

func (s *RpcServer) Start() (err error) {
	// 启动监听
	_, port, _ := net.SplitHostPort(s.Addr) // 没必要指定ip来监听
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		// todo
		return err
	}

	go func() {
		s.Add(1)
		s.AcceptRoutine()
		s.Done()
	}()

	// 服务已经启动，注册到注册中心，调用注册
	err = s.Register(s)
	if err != nil {
		return err
	}

	// 调用心跳
	go func() {
		s.Heartbeat(s)
	}()

	return
}

func (s *RpcServer) Stop() (err error) {
	s.Cancel()

	s.listener.Close()
	//s.coordinator.Stop()

	s.Wait()
	fmt.Println("exit success")
	return
}

func (s *RpcServer) AcceptRoutine() (err error) {
	var tempDelay time.Duration
	for {

		conn, err := s.listener.Accept()
		// 直接仿照 net/http 包的异常处理
		if err != nil {
			select {
			case <-s.Ctx.Done():
				return err
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		go s.serve(conn)

	}
}
func (s *RpcServer) serve(conn net.Conn) {
	s.Add(1)            // 防止处理一半 stop 退出被关掉，让 stop() 等一下;上游已经关闭listener，不用担心一直add导致这里退不出去
	rpc.ServeConn(conn) // ServeConn 会去close 掉 conn
	s.Done()
}

func GetIpv4_192_168() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), "192.168.") == 0 {
				return sip
			}
		}
	}
	return ""
}
func GetIpv4_172() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), "172.") == 0 {
				return sip
			}
		}
	}
	return ""
}
func GetIpv4_10() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ip, _, err := net.ParseCIDR(addr.String()); err == nil {
			sip := ip.To4().String()
			if strings.Index(addr.String(), "10.") == 0 {
				return sip
			}
		}
	}
	return ""
}
