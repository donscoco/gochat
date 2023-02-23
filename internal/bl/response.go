package bl

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
)

// 存放 api 服务的json响应逻辑格式

/////////////////////////////////////////

//////   存放 http 响应的通用 模版 //////
/*
{
	"code": int
	"message": string
	"data": interface{}
}
*/
type ResponseCode int

const (
	SuccessCode ResponseCode = 200
	// 请求成功

	// 请求错误
	// 自定义反馈错误信息 400
	ErrMsgCode ResponseCode = 400
	// 退出登陆 （前端回到登陆页面） 401
	InternalErrorCode ResponseCode = 401 // 用户没有登陆
	// http请求方式有误 405
	InvalidRequestErrorCode ResponseCode = 405
	// 非法操作
	IllegalOperationErrorCode ResponseCode = 406

	// 系统错误
	// 服务器出了点小差，请稍后再试 500
	UnknowCode ResponseCode = 500
	// 服务器不支持当前请求所需要的某个功能 501
	NotFoundServerCode ResponseCode = 501

	// 内部处理 4000+处理错误 5000+系统级错误
	OffLineCode ResponseCode = 4001
	NilCode     ResponseCode = 4002

	InvalidErrorCode ResponseCode = 5001 // 解析参数错误
	RPCErrorCode     ResponseCode = 5002 // rpc 接口调用失败

)

type HTTPResponse struct {
	Code    ResponseCode `json:"code"`
	Msg     string       `json:"message"`
	Data    interface{}  `json:"data"`
	TraceId interface{}  `json:"trace_id"` // 用于追踪定位问题
	//Stack     interface{}  `json:"stack"`
}

func ResponseError(c *gin.Context, code ResponseCode, err error) {

	// 放统一的响应处理逻辑

	trace, _ := c.Get("trace")
	if err == nil {
		err = errors.New("unknow err")
	}
	resp := &HTTPResponse{Code: code, Msg: err.Error(), Data: "", TraceId: trace}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
	c.AbortWithError(200, err)
}

func ResponseSuccess(c *gin.Context, data interface{}) {

	trace, _ := c.Get("trace")
	resp := &HTTPResponse{Code: SuccessCode, Msg: "成功", Data: data, TraceId: trace}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response))
}

func ResponseMsg(c *gin.Context, code ResponseCode, Msg string, data interface{}) {

	trace, _ := c.Get("trace")
	resp := &HTTPResponse{Code: code, Msg: Msg, Data: data, TraceId: trace}
	c.JSON(200, resp)
	response, _ := json.Marshal(resp)
	c.Set("response", string(response)) //用于给 log_req 记录
}

/////////////////////////////////////////
