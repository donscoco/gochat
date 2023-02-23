package dao

import (
	"github.com/donscoco/gochat/internal/bl"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type GroupMember struct {
	Id        int    `json:"id" gorm:"primary_key" description:"自增主键"`
	GroupId   int    `json:"group_id" gorm:"column:group_id"`
	UserId    int    `json:"user_id" gorm:"column:user_id"`
	AliasName string `json:"alias_name" gorm:"column:alias_name"`
	Remark    string `json:"remark" gorm:"column:remark"`

	//UpdateTime     time.Time `json:"update_time" gorm:"column:update_time" description:"更新时间"`
	CreatedTime time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *GroupMember) TableName() string {
	return "chat_group_member"
}

func (t *GroupMember) Find(c *gin.Context, tx *gorm.DB, search *GroupMember) (*GroupMember, error) {
	out := &GroupMember{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 根据where 查找
func (t *GroupMember) FindByWhere(c *gin.Context, tx *gorm.DB, query interface{}, args ...interface{}) (*GroupMember, error) {
	out := &GroupMember{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Where(query, args...).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 根据where 查找
func (t *GroupMember) SelectByWhere(c *gin.Context, tx *gorm.DB, query interface{}, args ...interface{}) (*GroupMember, error) {
	out := &GroupMember{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where(query, args...).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 插入
func (t *GroupMember) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Save(t).Error
}

func (t *GroupMember) CreateInBatches(c *gin.Context, tx *gorm.DB, members []GroupMember) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).
		CreateInBatches(members, len(members)).Error
}

// https://gorm.io/zh_CN/docs/update.html
func (t *GroupMember) Update(c *gin.Context, tx *gorm.DB, col string, val interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Update(col, val).Error
}

// 通过id 来作为 where 条件
func (t *GroupMember) UpdateById(c *gin.Context, tx *gorm.DB, m map[string]interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Updates(m).Error // where 只会设置主键id。
}
func (t *GroupMember) UpdateByWhere(c *gin.Context, tx *gorm.DB, m map[string]interface{}, query interface{}, args ...interface{}) error {
	// 更新类型不要创建新的session instance，因为外面可能是启动一个事务来处理。
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Where(query, args...).Updates(m).Error // where 只会设置主键id。
}

func (t *GroupMember) ListGroupMemberByUserId(c *gin.Context, tx *gorm.DB, uid int) ([]GroupMember, error) {

	var list []GroupMember

	query := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where("user_id=? AND is_delete=?", uid, 0)

	err := query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

func (t *GroupMember) ListGroupMemberByGroupId(c *gin.Context, tx *gorm.DB, gid int) ([]GroupMember, error) {

	var list []GroupMember

	query := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where("group_id=? AND is_delete=?", gid, 0)

	err := query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}
