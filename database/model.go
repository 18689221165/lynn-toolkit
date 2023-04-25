package database

import (
	"gorm.io/gorm"
	"youliao.cn/liaoma-toolkit/types"
)

// Model a basic GoLang struct which includes the following fields: ID, CreatedAt, UpdatedAt, DeletedAt
// It may be embedded into your model or you may build your own model without it
//    type User struct {
//      gorm.Model
//    }
type Model struct {
	ID        uint64     `gorm:"primaryKey"`
	CreatedAt types.Time `gorm:"index:idx_ctime"`
	UpdatedAt types.Time
}

func (o *Model) BeforeCreate(tx *gorm.DB) error {
	o.CreatedAt = types.NowTime()
	return nil
}

func (o *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = types.NowTime()
	return nil
}
