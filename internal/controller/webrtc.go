package controller

import (
	"encoding/json"
	"fmt"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/internal/module/ws_sender"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type WebRTCController struct{}

func WebRTCRegister(group *gin.RouterGroup) {
	webRTCController := &WebRTCController{}
	group.GET("/private/iceservers", webRTCController.GetICEServers) // 查询 stun 服务器

	group.POST("/private/call", webRTCController.Call)
	group.POST("/private/cancel", webRTCController.Cancel)
	group.POST("/private/reject", webRTCController.Reject)
	group.POST("/private/accept", webRTCController.Accept)
	group.POST("/private/candidate", webRTCController.Candidate)
	group.POST("/private/failed", webRTCController.Failed)
	group.POST("/private/handup", webRTCController.Handup)
}

func (ice *WebRTCController) GetICEServers(c *gin.Context) {
	// todo 搭建自己的stun
	data := []map[string]string{
		map[string]string{
			"urls":       "stun:www.boxim.online:3478",
			"username":   "admin",
			"credential": "admin123",
		},
		map[string]string{
			"urls":       "turn:www.boxim.online:3478",
			"username":   "admin",
			"credential": "admin123",
		},
		//{
		//	"username":   "1675883909:sample",
		//	"credential": "OBtiWJQ4uuU20VCwBsS5nKgQWJ4",
		//	"urls":       "turn:127.0.0.1:19302",
		//},
	}

	bl.ResponseSuccess(c, data)
}

/*
http://localhost/api/webrtc/private/call?uid=67
Request Method: POST
{"sdp":"","type":"offer"}
{"code":200,"message":"成功","data":null}
*/
func (ice *WebRTCController) Call(c *gin.Context) {
	// 获取目标uid，发送websocket信息，

	// 获得请求内容
	uidstr := c.Query("uid")
	userId, err := strconv.Atoi(uidstr)
	if err != nil {
		//
		bl.ResponseError(c, 2001, err)
		return
	}

	params := &model.WebRTCOffer{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	ch, ok := ws_sender.GetChan(uidstr)
	if !ok {

		//不在线就发送websocket 给自己
		selfCh, ok := ws_sender.GetChan(strconv.Itoa(sessionInfo.ID))
		if !ok { // 这里肯定有的，没有就是出现bug了
			bl.ResponseMsg(c, bl.InvalidErrorCode, "系统错误", nil)
		}

		t := time.Now().Unix()
		webRTCOutput := model.WebRTCOfferOutput{
			Content:  "对方当前不在线",
			Id:       0,
			RecvId:   sessionInfo.ID,
			SendId:   userId,
			SendTime: &t,
			Type:     model.RTC_FAILED, // 对应前端的code
		}
		webrtcMsg := model.WebSocketResponse{
			Cmd:  3, // todo 看下前端什么意思
			Data: webRTCOutput,
		}
		msg, _ := json.Marshal(webrtcMsg)

		selfCh <- msg
		bl.ResponseSuccess(c, nil)
		return
	}

	content, _ := json.Marshal(params)

	webRTCOutput := model.WebRTCOfferOutput{
		Content:  string(content),
		Id:       0,
		RecvId:   userId,
		SendId:   sessionInfo.ID,
		SendTime: nil,
		Type:     model.RTC_CALL, // 对应前端的code
	}
	webrtcMsg := model.WebSocketResponse{
		Cmd:  3, // todo 看下前端什么意思
		Data: webRTCOutput,
	}

	msg, _ := json.Marshal(webrtcMsg)

	ch <- msg

	bl.ResponseSuccess(c, nil)
}

/*
主动取消
Request URL: http://localhost/api/webrtc/private/cancel?uid=67
Request Method: POST
{"code":200,"message":"成功","data":null}
*/
func (ice *WebRTCController) Cancel(c *gin.Context) {
	// 获得请求内容
	uidstr := c.Query("uid")
	userId, err := strconv.Atoi(uidstr)
	if err != nil {
		// todo
		bl.ResponseError(c, 2001, err)
		return
	}

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	ch, ok := ws_sender.GetChan(uidstr)
	if !ok {
		// 不在线，他取消也不管，一般在不在线的call前端就直接取消了。前端逻辑能走到这里肯定是在线的。但是，如果有人故意调用这个接口，我们直接返回正常就行
		bl.ResponseSuccess(c, nil)
		return
	}

	webRTCOutput := model.WebRTCOfferOutput{
		Content:  "",
		Id:       0,
		RecvId:   userId,
		SendId:   sessionInfo.ID,
		SendTime: nil,
		Type:     model.RTC_CANCEL, // 对应前端的code
	}
	webrtcMsg := model.WebSocketResponse{
		Cmd:  3, // todo 看下前端什么意思
		Data: webRTCOutput,
	}

	msg, _ := json.Marshal(webrtcMsg)

	ch <- msg
	bl.ResponseSuccess(c, nil)
	return
}

/*
Request URL: http://localhost/api/webrtc/private/reject?uid=68
Request Method: POST
*/
func (ice *WebRTCController) Reject(c *gin.Context) {
	// 获得请求内容
	uidstr := c.Query("uid")
	userId, err := strconv.Atoi(uidstr)
	if err != nil {
		// todo
		bl.ResponseError(c, 2001, err)
		return
	}

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	ch, ok := ws_sender.GetChan(uidstr)
	if !ok {
		// 不在线，他取消也不管，一般在不在线的call前端就直接取消了。前端逻辑能走到这里肯定是在线的。但是，如果有人故意调用这个接口，我们直接返回正常就行
		bl.ResponseSuccess(c, nil)
		return
	}

	webRTCOutput := model.WebRTCOfferOutput{
		Content:  "",
		Id:       0,
		RecvId:   userId,
		SendId:   sessionInfo.ID,
		SendTime: nil,
		Type:     model.RTC_REJECT, // 对应前端的code
	}
	webrtcMsg := model.WebSocketResponse{
		Cmd:  3, // todo 看下前端什么意思
		Data: webRTCOutput,
	}

	msg, _ := json.Marshal(webrtcMsg)

	ch <- msg
	bl.ResponseSuccess(c, nil)
	return
}

/*
Request URL: http://localhost/api/webrtc/private/accept?uid=68
Request Method: POST
{"sdp":"","type":"offer"}
{"code":200,"message":"成功","data":null}
*/
func (ice *WebRTCController) Accept(c *gin.Context) {
	// 获取目标uid，发送websocket信息，

	// 获得请求内容
	uidstr := c.Query("uid")
	userId, err := strconv.Atoi(uidstr)
	if err != nil {
		// todo
		bl.ResponseError(c, 2001, err)
		return
	}

	params := &model.WebRTCOffer{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	ch, ok := ws_sender.GetChan(uidstr)
	if !ok {

		//不在线就发送websocket 给自己
		selfCh, ok := ws_sender.GetChan(strconv.Itoa(sessionInfo.ID))
		if !ok { // 这里肯定有的，没有就是出现bug了
			bl.ResponseMsg(c, bl.InvalidErrorCode, "系统错误", nil)
		}

		t := time.Now().Unix()
		webRTCOutput := model.WebRTCOfferOutput{
			Content:  "对方当前不在线",
			Id:       0,
			RecvId:   sessionInfo.ID,
			SendId:   userId,
			SendTime: &t,
			Type:     model.RTC_FAILED, // 对应前端的code
		}
		webrtcMsg := model.WebSocketResponse{
			Cmd:  3, // todo 看下前端什么意思
			Data: webRTCOutput,
		}
		msg, _ := json.Marshal(webrtcMsg)

		selfCh <- msg
		bl.ResponseSuccess(c, nil)
		return
	}

	content, _ := json.Marshal(params)

	webRTCOutput := model.WebRTCOfferOutput{
		Content:  string(content),
		Id:       0,
		RecvId:   userId,
		SendId:   sessionInfo.ID,
		SendTime: nil,
		Type:     model.RTC_ACCEPT, // 对应前端的code
	}
	webrtcMsg := model.WebSocketResponse{
		Cmd:  3, // todo 看下前端什么意思
		Data: webRTCOutput,
	}

	msg, _ := json.Marshal(webrtcMsg)

	ch <- msg

	bl.ResponseSuccess(c, nil)
}

/*
Request URL: http://localhost/api/webrtc/private/candidate?uid=67
Request Method: POST

	{
		"candidate": string
		"sdpMid": string
		"sdpMLineIndex": int
	}

{"code":200,"message":"成功","data":null}
*/
func (ice *WebRTCController) Candidate(c *gin.Context) { // todo 还要再检查一下啊
	// 获取目标uid，发送websocket信息，

	// 获得请求内容
	uidstr := c.Query("uid")
	userId, err := strconv.Atoi(uidstr)
	if err != nil {
		//
		bl.ResponseError(c, 2001, err)
		return
	}

	params := &model.WebRTCCandidate{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	ch, ok := ws_sender.GetChan(uidstr)
	if !ok {
		//不在线就 退出
		bl.ResponseSuccess(c, nil)
		return
	}

	content, _ := json.Marshal(params)

	webRTCOutput := model.WebRTCCandidateOutput{
		Content:  string(content),
		Id:       0,
		RecvId:   userId,
		SendId:   sessionInfo.ID,
		SendTime: nil,
		Type:     model.RTC_CANDIDATE, // 对应前端的code
	}
	webrtcMsg := model.WebSocketResponse{
		Cmd:  3, // todo 看下前端什么意思
		Data: webRTCOutput,
	}

	msg, _ := json.Marshal(webrtcMsg)

	ch <- msg

	bl.ResponseSuccess(c, nil)
}

/*
http://localhost/api/webrtc/private/failed?uid=undefined&reason=%E5%AF%B9%E6%96%B9%E6%AD%A3%E5%BF%99,%E6%9A%82%E6%97%B6%E6%97%A0%E6%B3%95%E6%8E%A5%E5%90%AC
Request Method: POST
*/
func (ice *WebRTCController) Failed(c *gin.Context) {
	// todo 还要再检查一下
	// 获取目标uid，发送websocket信息，

	// 获得请求内容
	//uidstr := c.Query("uid")
	//userId, err := strconv.Atoi(uidstr)
	//if err != nil {
	//	// todo
	//	bl.ResponseError(c, 2001, err)
	//	return
	//}
	//reasonstr := c.Query("reason")

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// todo websocket 发送忙

	bl.ResponseMsg(c, bl.SuccessCode, "对方正忙,暂时无法接听", nil)
	return
	//bl.ResponseSuccess(c, nil)
}

/*
Request URL: http://localhost/api/webrtc/private/handup?uid=68
Request Method: POST

{"code":200,"message":"成功","data":null}
*/
func (ice *WebRTCController) Handup(c *gin.Context) {
	// 获得请求内容
	uidstr := c.Query("uid")
	userId, err := strconv.Atoi(uidstr)
	if err != nil {
		// todo
		bl.ResponseError(c, 2001, err)
		return
	}

	//获得请求的用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(bl.SessionInfoKey)
	sessionInfo := &model.SessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), sessionInfo); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	ch, ok := ws_sender.GetChan(uidstr)
	if !ok {
		// 不在线，他取消也不管，一般在不在线的call前端就直接取消了。前端逻辑能走到这里肯定是在线的。但是，如果有人故意调用这个接口，我们直接返回正常就行
		bl.ResponseSuccess(c, nil)
		return
	}

	webRTCOutput := model.WebRTCOfferOutput{
		Content:  "",
		Id:       0,
		RecvId:   userId,
		SendId:   sessionInfo.ID,
		SendTime: nil,
		Type:     model.RTC_HANDUP, // 对应前端的code
	}
	webrtcMsg := model.WebSocketResponse{
		Cmd:  3, // todo 看下前端什么意思
		Data: webRTCOutput,
	}

	msg, _ := json.Marshal(webrtcMsg)

	ch <- msg
	bl.ResponseSuccess(c, nil)
	return
}
