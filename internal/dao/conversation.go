package dao

import (
	"github.com/donscoco/gochat/internal/bl"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	Id int `json:"id" gorm:"primary_key" description:"自增主键"`

	UserId  int `json:"user_id" gorm:"column:user_id"`
	RecvId  int `json:"recv_id" gorm:"column:recv_id"`   // type=1 私聊人id，type=2 群聊id
	GroupId int `json:"group_id" gorm:"column:group_id"` // type=1 私聊人id，type=2 群聊id
	Type    int `json:"type" gorm:"column:type"`

	ReadOffset  int64 `json:"read_offset" gorm:"column:read_offset"`
	WriteOffset int64 `json:"write_offset" gorm:"column:write_offset"`

	CreatedTime time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Conversation) TableName() string {
	return "chat_conversation"
}

func (t *Conversation) Find(c *gin.Context, tx *gorm.DB, search *Conversation) (*Conversation, error) {
	out := &Conversation{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 根据where 查找
func (t *Conversation) FindByWhere(c *gin.Context, tx *gorm.DB, query interface{}, args ...interface{}) (*Conversation, error) {
	out := &Conversation{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Where(query, args...).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 根据where 查找
func (t *Conversation) SelectByWhere(c *gin.Context, tx *gorm.DB, query interface{}, args ...interface{}) (*Conversation, error) {
	out := &Conversation{}
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
func (t *Conversation) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Save(t).Error
}

func (t *Conversation) CreateInBatches(c *gin.Context, tx *gorm.DB, members []Conversation) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).
		CreateInBatches(members, len(members)).Error
}

// https://gorm.io/zh_CN/docs/update.html
func (t *Conversation) Update(c *gin.Context, tx *gorm.DB, col string, val interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Update(col, val).Error
}

// 通过id 来作为 where 条件
func (t *Conversation) UpdateById(c *gin.Context, tx *gorm.DB, m map[string]interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Updates(m).Error // where 只会设置主键id。
}
func (t *Conversation) UpdateByWhere(c *gin.Context, tx *gorm.DB, m map[string]interface{}, query interface{}, args ...interface{}) error {
	// 更新类型不要创建新的session instance，因为外面可能是启动一个事务来处理。
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Where(query, args...).Updates(m).Error // where 只会设置主键id。
}
