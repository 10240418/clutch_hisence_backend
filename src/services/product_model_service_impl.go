package services

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type ProductModelService struct {
	db *gorm.DB
}

func NewProductModelService(db *gorm.DB) (IProductModelService, error) {
	return &ProductModelService{db: db}, nil
}

func (s *ProductModelService) CreateProductModel(productModel *models.ProductModel) error {
	return s.db.Create(productModel).Error
}

func (s *ProductModelService) GetProductModel(id int64) (*models.ProductModel, error) {
	var productModel models.ProductModel
	err := s.db.First(&productModel, id).Error
	return &productModel, err
}

func (s *ProductModelService) GetProductModels(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductModel, models.PaginationResult, error) {
	var productModels []models.ProductModel
	var pagination models.PaginationResult
	var model = s.db.Model(&models.ProductModel{}).Preload("Supplier")

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&productModels)
	if result.Error != nil {
		return []models.ProductModel{}, pagination, result.Error
	}

	return productModels, pagination, nil
}

func (s *ProductModelService) GetProductModelBySN(sn string) (*models.ProductModel, error) {
	var productModel models.ProductModel
	err := s.db.Where("sn = ?", sn).First(&productModel).Error
	return &productModel, err
}

func (s *ProductModelService) UpdateProductModel(productModelInstance *models.ProductModel, productModel map[string]interface{}) error {
	result := s.db.Model(productModelInstance).Updates(productModel)
	return result.Error
}

func (s *ProductModelService) DeleteProductModels(ids []int64) error {
	result := s.db.Delete(&models.ProductModel{}, ids)
	return result.Error
}
