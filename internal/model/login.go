package model

import (
	"time"
)

// 登陆结构的请求和响应

type SessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"user_name"`
	LoginTime time.Time `json:"login_time"`
}

type LoginInput struct {
	//UserName string `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //用户名

	UserName string `json:"username" form:"username" comment:"管理员用户名" example:"admin" validate:"required,valid_username"` //用户名
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                   //密码
}
type LoginOutput struct {
	// todo 暂时不输出什么
}

// todo 注册
// 请求结构体
/*
{
    "userName":"ironhead",
    "nickName":"铁头班小朋友",
    "password":"123456",
    "confirmPassword":"123456"
}
{"code":200,"message":"成功","data":null}
{"code":10004,"message":"该用户名已注册","data":null}
*/
// 注册
type RegisterInput struct {
	UserName        string `json:"userName" form:"userName" comment:"用户名" example:"admin" validate:"required,valid_username"` //用户名
	NickName        string `json:"nickName" form:"nickName" comment:"昵称" example:"admin" validate:"required,valid_username"`  //用户名
	Password        string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`                //密码
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" comment:"密码" example:"123456" validate:"required"`  //密码
}
