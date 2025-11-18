package models

// API 对应 'API' 表
type API struct {
	ModelFields `s2m:"-"`
	Name        string `gorm:"type:char(64)" json:"name"`
	AppID       string `gorm:"type:char(32)" json:"appId"`
	Secret      string `json:"secret,omitempty"`
}
