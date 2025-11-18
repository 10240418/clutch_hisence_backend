package utils

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"gorm.io/gorm"
)

func DoPagination(model *gorm.DB, paginate map[string]interface{}) (*gorm.DB, models.PaginationResult) {
	var total int64
	model.Count(&total)

	var pageSize int
	var pageNum int
	if value, ok := paginate["page_size"]; ok {
		pageSize = value.(int)
		model = model.Limit(pageSize)

		if value, ok := paginate["page_num"]; ok {
			pageNum = value.(int)
			model = model.Offset((pageNum - 1) * pageSize)
		}
	}

	pagination := models.PaginationResult{
		Total:    int(total),
		PageSize: pageSize,
		PageNum:  pageNum,
	}

	return model, pagination
}

func DoOrder(model *gorm.DB, paginate map[string]interface{}) *gorm.DB {
	var order string = "created_at DESC"
	if value, ok := paginate["asc"]; ok {
		if value.(bool) {
			order = "created_at ASC"
		} else {
			order = "created_at DESC"
		}
	}
	model = model.Order(order)

	return model
}
