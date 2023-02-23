package controller

import (
	"encoding/json"
	"time"

	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/dao"
	"github.com/donscoco/gochat/internal/model"
	"github.com/donscoco/gochat/pkg/gorm"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginController struct{}

func LoginRegister(group *gin.RouterGroup) {
	adminLogin := &LoginController{}
	group.POST("/login", adminLogin.Login)
	group.GET("/logout", adminLogin.LoginOut)
	group.POST("/register", adminLogin.Register)
}

func (userlogin *LoginController) Login(c *gin.Context) {

	// 获取body输入对象，db检查登陆，设置session，返回响应
	params := &model.LoginInput{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}

	user := &dao.User{}
	user, err = user.LoginCheck(c, tx, params)
	if err != nil {
		bl.ResponseError(c, 2002, err)
		return
	}

	//设置session
	sessInfo := &model.SessionInfo{
		ID:        user.Id,
		UserName:  user.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		bl.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c) // 获取cookie
	sess.Set(bl.SessionInfoKey, string(sessBts))
	sess.Save()

	out := &model.LoginOutput{
		// 暂时没有什么输出的
	}

	bl.ResponseSuccess(c, out)
}

func (userlogin *LoginController) LoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(bl.SessionInfoKey) // 删除在redis上的session
	sess.Save()
	bl.ResponseSuccess(c, "")
}

func (lc *LoginController) Register(c *gin.Context) {

	params := &model.RegisterInput{}
	if err := bl.DefaultGetValidParams(c, params); err != nil {
		bl.ResponseError(c, 2000, err)
		return
	}

	// 获取连接池的一个连接
	tx, err := gorm.GetGormPool("default")
	if err != nil {
		bl.ResponseError(c, 2001, err)
		return
	}

	// 1.检查用户名，2.写入
	search := &dao.User{}
	user, err := search.Register(c, tx, params)
	if err != nil {
		bl.ResponseError(c, 2002, err)
		return
	}

	//out := &model.LoginOutput{
	//	// 暂时没有什么输出的
	//}
	out := user // for debug

	bl.ResponseSuccess(c, out)
}
