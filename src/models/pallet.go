package models

// Pallet 对应 'Pallet' 表
type Pallet struct {
	ModelFields    `s2m:"-"`
	SN             string        `gorm:"type:char(32)" json:"sn"`
	ProductModelID *uint         `json:"productModelId"`
	ProductModel   *ProductModel `gorm:"foreignKey:ProductModelID" json:"productModel" s2m:"-"`
	ProductLineID  *uint         `json:"productLineId"`
	ProductLine    *ProductLine  `gorm:"foreignKey:ProductLineID" json:"productLine" s2m:"-"`
	Goal           int           `json:"goal"`
}
