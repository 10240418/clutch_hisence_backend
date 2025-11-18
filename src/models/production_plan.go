package models

import "time"

// ProductionPlan 对应 'ProductionPlan' 表
type ProductionPlan struct {
	ModelFields    `s2m:"-"`
	StartAt        time.Time    `gorm:"type:date" json:"startAt"`
	EndAt          time.Time    `gorm:"type:date" json:"endAt"`
	BelongsTo      string       `gorm:"type:char(10)" json:"belongsTo"`
	ProductModelID uint         `json:"productModelId"`
	ProductModel   ProductModel `gorm:"foreignKey:ProductModelID" json:"productModel" s2m:"-"`
	Planned        int          `json:"planned"`
	Actual         int          `gorm:"-" json:"actual" s2m:"-"`
}
