package dao

import (
	"github.com/donscoco/gochat/internal/bl"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Contacts struct {
	Id             int    `json:"id" gorm:"primary_key" description:"自增主键"`
	UserId         int    `json:"user_id" gorm:"column:user_id" description:"自增主键"`
	FriendId       int    `json:"friend_id" gorm:"column:friend_id" description:"自增主键"`
	FriendNickName string `json:"friend_nick_name" gorm:"column:friend_nick_name" description:"备注昵称"`
	Type           int    `json:"type" gorm:"column:type" description:"关系类型"`

	//UpdatedAt     time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedTime time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Contacts) TableName() string {
	return "chat_contacts"
}

func (t *Contacts) Find(c *gin.Context, tx *gorm.DB, search *Contacts) (*Contacts, error) {
	out := &Contacts{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Contacts) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Save(t).Error
}

func (t *Contacts) ListContactsByUser(tx *gorm.DB, uid int) ([]Contacts, error) {

	var list []Contacts

	query := tx.
		//Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where("user_id=? AND is_delete=?", uid, 0)

	err := query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

func (t *Contacts) FindByWhere(c *gin.Context, tx *gorm.DB, query interface{}, args ...interface{}) (*Contacts, error) {
	out := &Contacts{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).Where(query, args...).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Contacts) UpdateByWhere(c *gin.Context, tx *gorm.DB, m map[string]interface{}, query interface{}, args ...interface{}) error {
	// 更新类型不要创建新的session instance，因为外面可能是启动一个事务来处理。
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Where(query, args...).Updates(m).Error // where 只会设置主键id。
}

func (t *Contacts) DeleteContacts(c *gin.Context, tx *gorm.DB, m map[string]interface{}, userId, friendId int) error {
	ids := []int{userId, friendId}
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Where("user_id in ? AND friend_id in ? ", ids, ids).Updates(m).Error
}

func (t *Contacts) CreateInBatches(c *gin.Context, tx *gorm.DB, members []Contacts) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).CreateInBatches(members, len(members)).Error
}
