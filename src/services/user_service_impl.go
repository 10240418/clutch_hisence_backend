package services

import (
	"errors"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) (IUserService, error) {
	return &UserService{db: db}, nil
}

func (s *UserService) CreateUser(user *models.User) error {
	// Check if user exists
	result := s.db.Where("email = ?", user.Email).First(user)
	if result.Error == nil {
		return errors.New("user already exists")
	}

	if err := s.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

type UserIdentifierType int

const (
	IdentifierTypeID UserIdentifierType = iota
	IdentifierTypeEmail
	IdentifierTypeMobile
)

func (s *UserService) GetUserBy(identifierType UserIdentifierType, value interface{}) (*models.User, error) {
	var user models.User
	query := s.db.Model(&models.User{})

	switch identifierType {
	case IdentifierTypeID:
		query = query.Where("id = ?", value)
	case IdentifierTypeEmail:
		query = query.Where("email = ?", value)
	case IdentifierTypeMobile:
		query = query.Where("mobile = ?", value)
	default:
		return nil, errors.New("invalid identifier type")
	}

	result := query.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (s *UserService) GetUsers(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.User, models.PaginationResult, error) {
	var users []models.User = []models.User{}
	var paginateResult models.PaginationResult
	model := s.db.Model(&models.User{})

	// Search
	if value, ok := query["keyword"]; ok {
		model = model.Where("email LIKE ? OR username LIKE ?", "%"+value.(string)+"%", "%"+value.(string)+"%")
	}
	delete(query, "keyword")
	model = model.Where(query)

	model, paginateResult = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&users)
	if result.Error != nil {
		return nil, models.PaginationResult{}, result.Error
	}

	return users, paginateResult, nil
}

func (s *UserService) UpdateUser(userObj *models.User, user map[string]interface{}) error {
	return s.db.Model(userObj).Updates(user).Error
}

func (s *UserService) DeleteUsers(ids []int64) error {
	result := s.db.Delete(&models.User{}, ids)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
