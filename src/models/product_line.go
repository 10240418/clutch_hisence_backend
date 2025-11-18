package models

// ProductLine 对应 'ProductLine' 表
type ProductLine struct {
	ModelFields    `s2m:"-"`
	Name           string `gorm:"type:char(64)" json:"name,omitempty"`
	PalletSnPrefix string `gorm:"type:char(16)" json:"palletSnPrefix,omitempty"`
	DeviceID       string `gorm:"type:char(64);unique" json:"deviceId"`
	IsRegistered   bool   `gorm:"default:false" json:"isRegistered"`
	PublicKey      string `gorm:"type:text" json:"publicKey,omitempty"`
}
