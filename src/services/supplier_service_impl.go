package services

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type SupplierService struct {
	db *gorm.DB
}

func NewSupplierService(db *gorm.DB) (ISupplierService, error) {
	return &SupplierService{db: db}, nil
}

func (s *SupplierService) CreateSupplier(supplier *models.Supplier) error {
	return s.db.Create(supplier).Error
}

func (s *SupplierService) GetSupplier(id int64) (*models.Supplier, error) {
	var supplier models.Supplier
	err := s.db.First(&supplier, id).Error
	return &supplier, err
}

func (s *SupplierService) GetSuppliers(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.Supplier, models.PaginationResult, error) {
	var suppliers []models.Supplier
	var pagination models.PaginationResult
	var model = s.db.Model(&models.Supplier{})

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&suppliers)
	if result.Error != nil {
		return []models.Supplier{}, pagination, result.Error
	}

	return suppliers, pagination, nil
}

func (s *SupplierService) UpdateSupplier(supplierInstance *models.Supplier, supplier map[string]interface{}) error {
	result := s.db.Model(supplierInstance).Updates(supplier)
	return result.Error
}

func (s *SupplierService) DeleteSuppliers(ids []int64) error {
	result := s.db.Delete(&models.Supplier{}, ids)
	return result.Error
}
