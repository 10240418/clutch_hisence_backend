package services

import (
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"gorm.io/gorm"
)

type ProductionPlanService struct {
	db *gorm.DB
}

func NewProductionPlanService(db *gorm.DB) (IProductionPlanService, error) {
	return &ProductionPlanService{db: db}, nil
}

func (s *ProductionPlanService) CreateProductionPlan(productionPlan *models.ProductionPlan) error {
	return s.db.Create(productionPlan).Error
}

func (s *ProductionPlanService) GetProductionPlan(id int64) (*models.ProductionPlan, error) {
	var productionPlan models.ProductionPlan
	err := s.db.First(&productionPlan, id).Error
	return &productionPlan, err
}

func (s *ProductionPlanService) GetProductionPlans(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductionPlan, models.PaginationResult, error) {
	var productionPlans []models.ProductionPlan
	var pagination models.PaginationResult
	var model = s.db.Model(&models.ProductionPlan{}).Preload("ProductModel")

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)
	model = utils.DoOrder(model, paginate)

	result := model.Find(&productionPlans)
	if result.Error != nil {
		return []models.ProductionPlan{}, pagination, result.Error
	}

	// TODO: 用SQL查询进行统计
	// 为每个生产计划计算实际完成数量
	for i := range productionPlans {
		var count int64
		s.db.Model(&models.Product{}).Where("production_plan_id = ?", productionPlans[i].ID).Count(&count)
		productionPlans[i].Actual = int(count)
	}

	return productionPlans, pagination, nil
}

func (s *ProductionPlanService) UpdateProductionPlan(productionPlanInstance *models.ProductionPlan, productionPlan map[string]interface{}) error {
	result := s.db.Model(productionPlanInstance).Updates(productionPlan)
	return result.Error
}

func (s *ProductionPlanService) DeleteProductionPlans(ids []int64) error {
	result := s.db.Delete(&models.ProductionPlan{}, ids)
	return result.Error
}

func (s *ProductionPlanService) GetProductionPlansByDateRange(baseDate time.Time) (map[string][]models.ProductionPlan, error) {
	// 计算T, T+1, T+2, T+3的日期
	t := baseDate
	tPlus1 := baseDate.AddDate(0, 0, 1)
	tPlus2 := baseDate.AddDate(0, 0, 2)
	tPlus3 := baseDate.AddDate(0, 0, 3)

	// 查询指定日期范围内的生产计划
	var productionPlans []models.ProductionPlan
	err := s.db.Preload("ProductModel").
		Where("(start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?) OR (start_at <= ? AND end_at >= ?)",
			t.Format("2006-01-02"), t.Format("2006-01-02"), tPlus1.Format("2006-01-02"), tPlus1.Format("2006-01-02"), tPlus2.Format("2006-01-02"), tPlus2.Format("2006-01-02"), tPlus3.Format("2006-01-02"), tPlus3.Format("2006-01-02")).Debug().
		Find(&productionPlans).Error

	if err != nil {
		return nil, err
	}

	// 按日期分组生产计划
	result := map[string][]models.ProductionPlan{
		"T":   {},
		"T+1": {},
		"T+2": {},
		"T+3": {},
	}

	for _, plan := range productionPlans {
		// 检查计划是否包含T日期
		if plan.StartAt.Before(t.AddDate(0, 0, 1)) && plan.EndAt.After(t.AddDate(0, 0, -1)) {
			result["T"] = append(result["T"], plan)
		}
		// 检查计划是否包含T+1日期
		if plan.StartAt.Before(tPlus1.AddDate(0, 0, 1)) && plan.EndAt.After(tPlus1.AddDate(0, 0, -1)) {
			result["T+1"] = append(result["T+1"], plan)
		}
		// 检查计划是否包含T+2日期
		if plan.StartAt.Before(tPlus2.AddDate(0, 0, 1)) && plan.EndAt.After(tPlus2.AddDate(0, 0, -1)) {
			result["T+2"] = append(result["T+2"], plan)
		}
		// 检查计划是否包含T+3日期
		if plan.StartAt.Before(tPlus3.AddDate(0, 0, 1)) && plan.EndAt.After(tPlus3.AddDate(0, 0, -1)) {
			result["T+3"] = append(result["T+3"], plan)
		}
	}

	return result, nil
}

func (s *ProductionPlanService) GetActiveProductionPlan(date time.Time, productModelID *uint, allowExceed bool) (*models.ProductionPlan, error) {
	var productionPlan models.ProductionPlan

	if productModelID == nil {
		return nil, gorm.ErrRecordNotFound
	}

	query := s.db.Preload("ProductModel").
		Where("product_model_id = ? AND start_at <= ? AND end_at >= ?", productModelID, date.Format("2006-01-02"), date.Format("2006-01-02"))

	// 如果不允许超过数量限制，则添加数量检查条件
	if !allowExceed {
		query = query.Where("planned > (SELECT COUNT(*) FROM products WHERE production_plan_id = production_plans.id)")
	}

	err := query.First(&productionPlan).Error

	return &productionPlan, err
}
