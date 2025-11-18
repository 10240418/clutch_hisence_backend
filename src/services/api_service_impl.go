package services

import (
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type APIService struct {
	db *gorm.DB
}

func NewAPIService(db *gorm.DB) (IAPIService, error) {
	return &APIService{db: db}, nil
}

func (s *APIService) CreateAPI(api *models.API) error {
	return s.db.Create(api).Error
}

func (s *APIService) GetAPI(id int64) (*models.API, error) {
	var api models.API
	err := s.db.First(&api, id).Error
	return &api, err
}

func (s *APIService) GetAPIs(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.API, models.PaginationResult, error) {
	var apis []models.API
	var pagination models.PaginationResult
	var model = s.db.Model(&models.API{})

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&apis)
	if result.Error != nil {
		return []models.API{}, pagination, result.Error
	}

	return apis, pagination, nil
}

func (s *APIService) UpdateAPI(apiInstance *models.API, api map[string]interface{}) error {
	result := s.db.Model(apiInstance).Updates(api)
	return result.Error
}

func (s *APIService) DeleteAPIs(ids []int64) error {
	result := s.db.Delete(&models.API{}, ids)
	return result.Error
}
