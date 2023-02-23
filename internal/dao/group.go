package dao

import (
	"github.com/donscoco/gochat/internal/bl"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	Id        int    `json:"id" gorm:"primary_key" description:"自增主键"`
	OwnerId   int    `json:"owner_id" gorm:"column:owner_id"`
	GroupName string `json:"group_name" gorm:"column:group_name"`
	HeadImage string `json:"head_image" gorm:"column:head_image"`
	Notice    string `json:"notice" gorm:"column:notice"`

	//UpdateTime     time.Time `json:"update_time" gorm:"column:update_time" description:"更新时间"`
	CreatedTime time.Time `json:"create_time" gorm:"column:create_time" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Group) TableName() string {
	return "chat_group"
}

// 根据id 查找
func (t *Group) Find(c *gin.Context, tx *gorm.DB, search *Group) (*Group, error) {
	out := &Group{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 根据where 查找
func (t *Group) FindByWhere(c *gin.Context, tx *gorm.DB, query interface{}, args ...interface{}) (*Group, error) {
	out := &Group{}
	err := tx.Set("trace_context", bl.GetGinTraceContext(c)).Where(query, args...).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// 单条插入
func (t *Group) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Save(t).Error
}

// 批量插入
func (t *Group) CreateInBatches(c *gin.Context, tx *gorm.DB, members []Group) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).
		CreateInBatches(members, len(members)).Error
}

// 根据id 单列更新 （t 要求是已经设置的主键id的）
func (t *Group) Update(c *gin.Context, tx *gorm.DB, col string, val interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Update(col, val).Error
}

// 根据id 多列更新 （t 要求是已经设置的主键id的）
func (t *Group) UpdateById(c *gin.Context, tx *gorm.DB, m map[string]interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Updates(m).Error
}

// 根据where条件 多列更新
func (t *Group) UpdateByWhere(c *gin.Context, tx *gorm.DB, m map[string]interface{}, where interface{}, args ...interface{}) error {
	return tx.Set("trace_context", bl.GetGinTraceContext(c)).Model(t).Where(where, args...).Updates(m).Error
}

// 根据条件批量查找
func (t *Group) ListGroupByWhere(c *gin.Context, tx *gorm.DB, where interface{}, args ...interface{}) ([]Group, error) {
	var list []Group
	query := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where(where, args...)
	err := query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}

////////////// 上面是通用的，下面是为了方便业务逻辑特殊的 //////////////

// 场景：用户查询自己所属的所有群组
func (t *Group) ListGroupByOwnerId(c *gin.Context, tx *gorm.DB, uid int) ([]Group, error) {
	var list []Group
	query := tx.Set("trace_context", bl.GetGinTraceContext(c)).
		Table(t.TableName()).
		Select("*").
		Where("owner_id=? AND is_delete=?", uid, 0)
	err := query.Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return list, nil
}
