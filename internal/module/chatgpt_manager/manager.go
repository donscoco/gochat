package chatgpt_manager

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/donscoco/gochat/internal/module/conn_manager"
	"log"
	"sync"
	"time"
)

type ChatGPTClientManager struct {
	apikey string

	// map 用来存放调用的client
	Pool map[int]*gpt3.Client // 用来存放每个用户向chatgpt发起的消息

	//input    chan interface{}
	//callback func(resp *gpt3.CompletionResponse)

	Ctx    context.Context
	Cancel context.CancelFunc
	m      sync.Mutex
	wg     sync.WaitGroup
}

func CreateChatGPTClientManager(apikey string) (server *ChatGPTClientManager) {
	server = new(ChatGPTClientManager)

	server.Pool = make(map[int]*gpt3.Client)
	server.apikey = apikey
	//server.input = make(chan interface{})
	//server.callback = callback

	server.Ctx, server.Cancel = context.WithCancel(context.TODO())

	return

}

func (s *ChatGPTClientManager) Start() (err error) {

	// 守护协程 负责存活检测等处理
	go func() {
		s.wg.Add(1)
		s.Daemon()
		s.wg.Done()
	}()

	// todo log

	return

}

// tip 后续尝试另一种风格的安全退出，使用外部传入的context去做done的安全退出
func (s *ChatGPTClientManager) Stop() (err error) {
	s.Cancel()

	// close others

	s.wg.Wait()
	fmt.Println("exit success")
	return
}

// 守护协程，用来定时检查用户的连接，超时关闭等一些后台处理
func (s *ChatGPTClientManager) Daemon() (err error) {

	ticker := time.NewTicker(60 * time.Second) // todo 具体多长时间检测，最好改成配置文件来设置
	for {

		s.m.Lock() // todo 优化，例如：如果map很大，每个client都要进行一次网络io往返时间的检查，锁住的时间太长。
		for uid, _ := range s.Pool {
			if conn_manager.IsAlive(uid) {
				continue
			}
			delete(s.Pool, uid) //
		}
		s.m.Unlock()

		<-ticker.C
	}
}
func (s *ChatGPTClientManager) Process(uid int, content string) (result string, err error) {

	err = (*s.GetClient(uid)).CompletionStreamWithEngine(
		s.Ctx,
		gpt3.TextDavinci003Engine,
		//gpt3.TextAda001Engine,
		gpt3.CompletionRequest{
			Prompt: []string{
				content,
			},
			MaxTokens:   gpt3.IntPtr(3000),
			Temperature: gpt3.Float32Ptr(0),
		},
		// 响应回调函数
		func(resp *gpt3.CompletionResponse) {
			//fmt.Print(resp.Choices[0].Text)
			result = result + resp.Choices[0].Text
		},
	)
	log.Printf("chatgpt resp:%s\n", result)
	if err != nil {
		return "", err
	}
	return result, nil

}
func (s *ChatGPTClientManager) GetClient(uid int) (client *gpt3.Client) {
	s.m.Lock()
	defer s.m.Unlock()

	cli := s.Pool[uid]
	if cli != nil {
		return cli
	}

	// 没有就创建新的
	newCli := gpt3.NewClient(s.apikey)
	s.Pool[uid] = &newCli
	return s.Pool[uid]
}
func (s *ChatGPTClientManager) RemoveClient(uid int) {
	s.m.Lock()
	defer s.m.Unlock()

	delete(s.Pool, uid)
}

/////
