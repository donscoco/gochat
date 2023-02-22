package dao

import (
	"errors"
	"github.com/donscoco/gochat/internal/bl"
	"github.com/donscoco/gochat/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id       int    `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName string `json:"user_name" gorm:"column:user_name" description:"用户名"`
	//Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string `json:"password" gorm:"column:password" description:"密码"`
	NickName  string `json:"nick_name" gorm:"column:nick_name" description:"昵称"`
	HeadImage string `json:"head_image" gorm:"column:head_image" description:"头像"`

	Sex       int    `json:"sex" gorm:"column:sex" description:"性别"`
	Signature string `json:"signature" gorm:"column:signature" description:"签名"`

	LastLoginTime time.Time `json:"last_login_time" gorm:"column:last_login_time" description:"最后时间"`
	UpdatedAt     time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedTime   time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
	IsDelete      int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *User) TableName() string {
	return "chat_user"
}

func (t *User) Find(c *gin.Context, tx *gorm.DB, search *User) (*User, error) {
	out := &User{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (t *User) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Save(t).Error
}

// 通过id 来作为 where 条件,要在user{}中设置id
func (t *User) UpdateById(c *gin.Context, tx *gorm.DB, m map[string]interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Updates(m).Error // where 只会设置主键id。
}
func (t *User) UpdateByWhere(c *gin.Context, tx *gorm.DB, m map[string]interface{}, where interface{}, args ...interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Where(where, args...).Updates(m).Error // where 只会设置主键id。
}

// 功能

// //1. params.UserName 取得管理员信息 admininfo
// //2. admininfo.salt + params.Password sha256 => saltPassword
// //3. saltPassword==admininfo.password
func (t *User) LoginCheck(c *gin.Context, tx *gorm.DB, param *model.LoginInput) (*User, error) {
	adminInfo, err := t.Find(c, tx, (&User{UserName: param.UserName, IsDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	// todo 密码加盐
	//saltPassword := tool.GenSaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != param.Password {
		return nil, errors.New("密码错误，请重新输入")
	}
	return adminInfo, nil
}

func (t *User) Register(c *gin.Context, tx *gorm.DB, param *model.RegisterInput) (*User, error) {
	// todo 校验，密码等。虽然前端有校验，但是服务也要校验

	adminInfo, err := t.Find(c, tx, &User{UserName: param.UserName, IsDelete: 0})
	if err != nil {
		return nil, errors.New("查询出错")
	}
	if adminInfo != nil && adminInfo.UserName == param.UserName {
		return nil, errors.New("用户已存在")
	}

	user := &User{
		UserName:      param.UserName,
		Password:      param.Password,
		NickName:      param.NickName,
		LastLoginTime: time.Now(),
		UpdatedAt:     time.Now(),
		CreatedTime:   time.Now(),
		IsDelete:      0,
	}
	user.Save(c, tx)

	return user, nil
}

func (t *User) FindByNikeName(c *gin.Context, tx *gorm.DB, search *User) (list []*User, err error) {

	list = make([]*User, 0, 16)

	query := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where(" nick_name like ? AND is_delete = ? ", search.NickName+"%", 0)

	err = query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

func (t *User) FindByIds(c *gin.Context, tx *gorm.DB, ids []int) (list []*User, err error) {

	list = make([]*User, 0, 16)
	//where := "("
	//for _, id := range ids {
	//	where = where + strconv.Itoa(id) + ","
	//}
	//where = where[:len(where)-1] + ")"

	query := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where(" id in ? AND is_delete = ? ", ids, 0)

	err = query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}
