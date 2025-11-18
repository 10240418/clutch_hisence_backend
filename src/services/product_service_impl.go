package services

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) (IProductService, error) {
	return &ProductService{db: db}, nil
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.db.Create(product).Error
}

func (s *ProductService) GetProduct(id int64) (models.Product, error) {
	var product models.Product
	err := s.db.First(&product, id).Error
	return product, err
}

func (s *ProductService) GetProducts(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.Product, models.PaginationResult, error) {
	var products []models.Product
	var pagination models.PaginationResult
	var model = s.db.Model(&models.Product{}).Preload("ProductModel").Preload("ProductLine").Preload("Pallet").Preload("ProductionPlan").Preload("ProductModel.Supplier")

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&products)
	if result.Error != nil {
		return []models.Product{}, pagination, result.Error
	}

	return products, pagination, nil
}

func (s *ProductService) UpdateProduct(productInstance *models.Product, product map[string]interface{}) error {
	result := s.db.Model(productInstance).Updates(product)
	return result.Error
}

func (s *ProductService) DeleteProducts(ids []int64) error {
	result := s.db.Delete(&models.Product{}, ids)
	return result.Error
}
