package chatgpt_service

import (
	"bytes"
	"encoding/json"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/chatgpt_manager"
	"github.com/donscoco/gochat/pkg/gorm"
	"github.com/donscoco/gochat/pkg/iron_log"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

// obj
// start     配置+创建初始化obj
// stop

var ChatGPTClientManager *chatgpt_manager.ChatGPTClientManager

// 创建连接
var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, //连接超时
		KeepAlive: 30 * time.Second, //探活时间
	}).DialContext,
	MaxIdleConns:          100,              //最大空闲连接
	IdleConnTimeout:       90 * time.Second, //空闲超时时间
	TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
	ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
}

// var Client *http.Client
// todo url路径 和 一些数值参数 改成文件配置
func Start() (err error) {

	// 0.初始化chatGPT的 client
	// 1.查询判断db，看是否创建聊天机器人
	// 2.http登陆，对于的session
	// 3.http keep alive websocket
	// 4.接收websocket的 msg 类型数据结构

	ChatGPTClientManager = chatgpt_manager.CreateChatGPTClientManager(os.Getenv("API_KEY"))
	err = ChatGPTClientManager.Start()
	if err != nil {
		return err
	}

	// 初始化 机器人的账号
	user := initChatGPTRobot()

	// 创建客户端
	client := initClient(user)

	// 拿到cookie之后，设置一下 开启 webserver 轮询
	//创建一个拨号器，也可以用默认的 websocket.DefaultDialer
	dialer := websocket.DefaultDialer //websocket.Dialer{}
	dialer.Jar = client.Jar           // 设置下cookie
	conn, _, err := dialer.Dial("ws://localhost:9980/im/ws", nil)
	if err != nil {
		return err
	}

	// 创建channel，用来和conn通信
	recvChan := make(chan []byte)
	sendChan := make(chan []byte)

	// 用来轮询获取 recv buffer 数据
	go func() {
		for {
			// 第一个返回值是 消息类型，这里我们默认就是text了。后续如果扩展需要注意下
			_, message, err := conn.ReadMessage()
			if err != nil {
				// 关闭下游
				close(recvChan)
				return
			}
			recvChan <- message
		}
	}()

	// 用来定时发送心跳
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			<-ticker.C

			message := []byte(`{"cmd":1,"data":{"userId":99}}`) // 1:心跳，2:下线通知，3:私聊，4:群聊
			sendChan <- message
		}
	}()

	// 用来处理消息
	go func() {
		defer conn.Close()

		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case msg, ok := <-recvChan:
				if !ok { // 关闭下游
					return
				}

				// todo 这里要做异步处理，这里先用go简单做后续考虑抽出来做成一个任务队列，chatGPT 处理是比较慢的，这里不能因为一个客户的处理就阻塞整个流程
				go func() {
					job(msg, client)
				}()

			case msg, ok := <-sendChan:
				if !ok { // 关闭下游
					return
				}
				err = conn.WriteMessage(websocket.TextMessage, msg) // 默认是 text类型的消息
				if err != nil {
					// todo 这里发送的错误可能是超时，考虑是直接断开连接还是重试？
					return
				}
			case <-ticker.C:
				// todo 做一些检查等后台处理
			}
		}
	}()

	return
}
func Stop() (err error) {
	if ChatGPTClientManager != nil {
		return ChatGPTClientManager.Stop()
	}
	return nil
}

// 在数据库中查找 chatgpt的user，没有就创建
func initChatGPTRobot() (user *dao.User) {
	chatRobotId := 99
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		return
	}

	tx = tx.Begin()
	user = new(dao.User)
	//err = tx.Where("where id = ?", chatRobotId).Find(&user).Error
	err = tx.Raw("SELECT * FROM chat_user "+
		"WHERE id = ? AND is_delete = 0 ", chatRobotId).Scan(user).Error
	if err != nil {
		tx.Rollback()
		log.Fatalf(err.Error())
	}
	if user.Id == 0 { // 没有就创建
		user.Id = chatRobotId
		user.UserName = "chatGPT 机器人"
		user.Password = "qwer123456"
		user.Signature = "基于chatGPT实现的机器人"
		user.NickName = "chatGPT 机器人"
		user.LastLoginTime = time.Now()
		user.CreatedTime = time.Now()
		user.UpdatedAt = time.Now()
		err := tx.Save(&user).Error
		if err != nil {
			tx.Rollback()
			log.Fatalf(err.Error())
		}
	}
	tx.Commit()
	return
}

// 创建客户端
func initClient(user *dao.User) (client *http.Client) {
	// 1.登陆获得cookie
	// 2.给client 设置cookie
	client = &http.Client{
		Timeout:   time.Second * 30, //请求超时时间
		Transport: transport,
	}
	// 请求数据
	form := url.Values{}
	form.Add("username", user.UserName)
	form.Add("password", user.Password)
	resp, err := client.PostForm("http://localhost:9980/api/login", form)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse("http://localhost:9980/api/login")
	jar.SetCookies(u, resp.Cookies())
	client.Jar = jar ////顺便也给client 设置下cookie，后面机器人回复消息要验证session

	return
}

// 任务处理函数
func job(msg []byte, client *http.Client) (err error) {
	// 校验
	// 发送给chatGPT
	// 将结果发送消息给conn

	wsresp := new(model.WebSocketResponse)
	err = json.Unmarshal(msg, wsresp)
	if err != nil {
		iron_log.Warnf("[chatgpt] unmarshal process error: %+v", err)
	}
	if wsresp.Cmd != model.CmdPrivate {
		// 心跳，下线通知和群聊消息 不做回应
		//log.Println("cmd 不是 3 不做回应")
		return
	}

	//wsData := model.SendMsg{
	//	Content:  req.Content,
	//	Seq:      req.Seq,
	//	RecvId:   req.RecvId,
	//	GroupId:  req.GroupId,
	//	SendId:   req.SendId,
	//	SendTime: req.SendTime,
	//	Type:     req.Type,
	//}
	recvMsg := model.SendMsg{}
	tmp, _ := json.Marshal(wsresp.Data)
	err = json.Unmarshal(tmp, &recvMsg)
	if err != nil {
		iron_log.Warnf("[chatgpt] unmarshal error: %+v", err)
		return
	}
	chatGPTResp, err := ChatGPTClientManager.Process(recvMsg.SendId, recvMsg.Content)
	if err != nil {
		iron_log.Warnf("[chatgpt] client process error: %+v", err)
		return
	}

	// {"content":"hello","type":0,"recvId":99}
	// 暂时只处理私聊信息
	respMsg := model.MsgSendInput{}
	respMsg.Content = chatGPTResp
	respMsg.Type = 0 // 普通文本消息
	respMsg.RecvId = recvMsg.SendId

	url := "http://localhost:9980/api/message/private/send"
	data, err := json.Marshal(respMsg)
	if err != nil {
		iron_log.Errorf("[chatgpt] marshal error: %+v", err)
		return
	}

	resp, err := client.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		iron_log.Errorf("[chatgpt] send msg error: %+v", err)
		return
	}
	defer resp.Body.Close()
	d, _ := io.ReadAll(resp.Body)

	iron_log.Debugf("[chatgpt] send msg resp: %s", string(d))
	return

}
