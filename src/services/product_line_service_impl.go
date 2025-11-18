package services

import (
	"fmt"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type ProductLineService struct {
	db *gorm.DB
}

func NewProductLineService(db *gorm.DB) (IProductLineService, error) {
	return &ProductLineService{db: db}, nil
}

func (s *ProductLineService) CreateProductLine(productLine *models.ProductLine) error {
	return s.db.Create(productLine).Error
}

func (s *ProductLineService) GetProductLine(id int64) (*models.ProductLine, error) {
	var productLine models.ProductLine
	err := s.db.Debug().First(&productLine, id).Error
	fmt.Println(productLine)
	return &productLine, err
}

func (s *ProductLineService) GetProductLines(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductLine, models.PaginationResult, error) {
	var productLines []models.ProductLine
	var pagination models.PaginationResult
	var model = s.db.Model(&models.ProductLine{})

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&productLines)
	if result.Error != nil {
		return []models.ProductLine{}, pagination, result.Error
	}

	return productLines, pagination, nil
}

func (s *ProductLineService) UpdateProductLine(productLineInstance *models.ProductLine, productLine map[string]interface{}) error {
	result := s.db.Model(productLineInstance).Updates(productLine)
	return result.Error
}
func (s *ProductLineService) DeleteProductLines(ids []int64) error {
	// Hard delete the product lines
	result := s.db.Unscoped().Delete(&models.ProductLine{}, ids)
	return result.Error
}

func (s *ProductLineService) GetProductLineByDeviceID(deviceID string) (*models.ProductLine, error) {
	var productLine models.ProductLine
	err := s.db.Where("device_id = ?", deviceID).First(&productLine).Error
	return &productLine, err
}
