package models

import (
	"time"

	"gorm.io/gorm"
)

type ModelFields struct {
	ID        int64          `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type PaginationQuery struct {
	PageNum  int  `form:"pageNum" json:"pageNum"`
	PageSize int  `form:"pageSize" json:"pageSize"`
	Asc      bool `form:"asc" json:"asc"`
}

type PaginationResult struct {
	Total    int `form:"total" json:"total"`
	PageNum  int `form:"pageNum" json:"pageNum"`
	PageSize int `form:"pageSize" json:"pageSize"`
}
