package models

// ProductModel 对应 'ProductModel' 表
type ProductModel struct {
	ModelFields `s2m:"-"`
	SN          string    `gorm:"type:char(16)" json:"sn"`
	PartNumber  string    `gorm:"type:char(32)" json:"partNumber"`
	Description string    `gorm:"type:char(128)" json:"description"`
	SupplierID  *uint     `json:"supplierId"`
	Supplier    *Supplier `gorm:"foreignKey:SupplierID" json:"supplier" s2m:"-"`
}
