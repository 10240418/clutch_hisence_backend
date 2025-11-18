package services

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type PalletService struct {
	db *gorm.DB
}

func NewPalletService(db *gorm.DB) (IPalletService, error) {
	return &PalletService{db: db}, nil
}

func (s *PalletService) CreatePallet(pallet *models.Pallet) error {
	return s.db.Create(pallet).Error
}

func (s *PalletService) GetPallet(id int64) (*models.Pallet, error) {
	var pallet models.Pallet
	err := s.db.First(&pallet, id).Error
	return &pallet, err
}

func (s *PalletService) GetPallets(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.Pallet, models.PaginationResult, error) {
	var pallets []models.Pallet
	var pagination models.PaginationResult
	var model = s.db.Model(&models.Pallet{}).Preload("ProductModel").Preload("ProductLine").Preload("ProductModel.Supplier")

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&pallets)
	if result.Error != nil {
		return []models.Pallet{}, pagination, result.Error
	}

	return pallets, pagination, nil
}

func (s *PalletService) UpdatePallet(palletInstance *models.Pallet, pallet map[string]interface{}) error {
	result := s.db.Model(palletInstance).Updates(pallet)
	return result.Error
}

func (s *PalletService) DeletePallets(ids []int64) error {
	result := s.db.Delete(&models.Pallet{}, ids)
	return result.Error
}
