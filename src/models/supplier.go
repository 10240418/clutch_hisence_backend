package models

// SupplierType 定义供应商类型的枚举
type SupplierType string

// Supplier 对应 'Supplier' 表
type Supplier struct {
	ModelFields `s2m:"-"`
	Name        string       `gorm:"type:char(64)" json:"name"`
	SAP         string       `gorm:"type:char(16)" json:"sap"`
	Type        SupplierType `gorm:"type:enum('直接供应','贸易商')" json:"type" s2m:"-"` // 假设的枚举值
}
