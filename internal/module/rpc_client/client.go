package rpc_client

import (
	"context"
	"sync"
)

type RpcClient struct {
	Retries      int // 重试次数
	TimeoutMs    int // 调用超时时间 ms
	UpdateClient func(*RpcClient) error
	Watcher      func(*RpcClient) error

	// 管理的代理
	ClientAgent map[string]*ClientAgent
	Balance     []*ClientAgent
	cursor      int // 均衡轮询的游标

	Context context.Context
	Cancel  context.CancelFunc

	sync.WaitGroup
	sync.Mutex
}

func CreateRpcClient(Retries int, TimeoutMs int,
	watcher func(*RpcClient) error,
	updateAddr func(*RpcClient) error,
) (c *RpcClient, err error) {
	c = new(RpcClient)
	if Retries != 0 {
		c.Retries = Retries
	} else {
		c.Retries = 1
	}
	if TimeoutMs != 0 {
		c.TimeoutMs = TimeoutMs
	}
	c.Watcher = watcher
	c.UpdateClient = updateAddr

	c.ClientAgent = make(map[string]*ClientAgent)
	c.Balance = make([]*ClientAgent, 0, 16)

	c.Context, c.Cancel = context.WithCancel(context.TODO())

	return c, nil
}
func (c *RpcClient) Start() (err error) {
	err = c.UpdateClient(c)
	if err != nil {
		return
	}

	go func() {
		c.Add(1)
		c.Watcher(c)
		c.Done()
	}()

	return nil
}
func (c *RpcClient) Stop() (err error) {
	c.Cancel()

	for _, ca := range c.ClientAgent {
		ca.Close()
	}

	c.Wait()
	return
}
func (c *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	var err error
	var ca *ClientAgent
	for i := 0; i < c.Retries; i++ {
		// 获得agent的一个client , 调用
		ca, err = c.GetClientAgent() // 注意使用redis保存conn的时候，要用addr 来获取
		if err != nil {
			continue
		}
		err = ca.Call(method, args, reply, int64(c.TimeoutMs))
		if err != nil {
			continue
		} else {
			break
		}
	}
	return err
}
func (c *RpcClient) GetClientAgent() (ca *ClientAgent, err error) {

	len := len(c.Balance)
	if len == 0 {
		return nil, ErrNoServers
	}

	c.Lock()
	defer c.Unlock()

	c.cursor = (c.cursor + 1) % len
	return c.Balance[c.cursor], nil
}

// 特殊处理
// 不使用负载均衡，直接指定ip进行调用，因为 需要用对应机器的连接的fd
func (c *RpcClient) CallSpecificAddr(method string, args interface{}, reply interface{}, addr string) (err error) {
	for i := 0; i < c.Retries; i++ {
		// 获得agent的一个client , 调用
		ca, err := c.GetSpecificClientAgent(addr)
		if err != nil {
			continue
		}
		err = ca.Call(method, args, reply, int64(c.TimeoutMs))
		if err != nil {
			continue
		} else {
			break
		}
	}
	return err
}
func (c *RpcClient) GetSpecificClientAgent(addr string) (ca *ClientAgent, err error) {

	len := len(c.ClientAgent)
	if len == 0 {
		return nil, ErrNoServers
	}

	c.Lock()
	defer c.Unlock()

	ca, ok := c.ClientAgent[addr]
	if !ok {
		return nil, ErrNoServers

	}
	return ca, nil
}
