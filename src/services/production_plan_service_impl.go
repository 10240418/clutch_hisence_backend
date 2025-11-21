package services

import (
	"fmt"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/utils"
	"github.com/xuri/excelize/v2"
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

func (s *ProductionPlanService) BatchCreateProductionPlans(plans []models.ProductionPlan) ([]models.ProductionPlan, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&plans).Error
	})
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (s *ProductionPlanService) GetProductionPlan(id int64) (*models.ProductionPlan, error) {
	var productionPlan models.ProductionPlan
	err := s.db.First(&productionPlan, id).Error
	return &productionPlan, err
}

func (s *ProductionPlanService) GetProductionPlans(query map[string]interface{}, paginate map[string]interface{}, sqlHandler ...func(*gorm.DB) *gorm.DB) ([]models.ProductionPlan, models.PaginationResult, error) {
	var productionPlans []models.ProductionPlan
	var pagination models.PaginationResult
	var model = s.db.Model(&models.ProductionPlan{})

	for _, handler := range sqlHandler {
		model = handler(model)
	}
	model = model.Where(query)

	model, pagination = utils.DoPagination(model, paginate)

	// Order by ID to maintain import order
	model = model.Order("id ASC")

	result := model.Find(&productionPlans)
	if result.Error != nil {
		return []models.ProductionPlan{}, pagination, result.Error
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

/*
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
*/

func (s *ProductionPlanService) GetProductionPlansByDateRange(baseDate time.Time) (map[string][]models.ProductionPlan, error) {
	// Temporary stub
	return nil, nil
}

func (s *ProductionPlanService) GetProductionPlansByDate(date time.Time) ([]models.ProductionPlan, error) {
	// 1. 查询指定日期的所有生产计划
	var plans []models.ProductionPlan
	dateStr := date.Format("2006-01-02")
	err := s.db.Where("DATE(plan_date) = ?", dateStr).Find(&plans).Error
	if err != nil {
		return nil, err
	}

	// 如果没有生产计划，直接返回
	if len(plans) == 0 {
		return plans, nil
	}

	// 2. 收集所有 part_number 用于 SQL IN 查询
	partNumbers := make([]string, 0, len(plans))
	for _, plan := range plans {
		partNumbers = append(partNumbers, plan.PartNumber)
	}

	// 3. 使用 SUBSTRING_INDEX 提取 description 中第一个 "/" 之前的部分，然后分组统计
	type CountResult struct {
		PartNumberPrefix string
		Count            int
	}

	var countResults []CountResult
	err = s.db.Table("products").
		Select("SUBSTRING_INDEX(product_models.description, '/', 1) as part_number_prefix, COUNT(*) as count").
		Joins("INNER JOIN product_models ON products.product_model_id = product_models.id").
		Where("DATE(products.created_at) = ?", dateStr).
		Where("SUBSTRING_INDEX(product_models.description, '/', 1) IN ?", partNumbers).
		Group("part_number_prefix").
		Scan(&countResults).Error

	if err != nil {
		return nil, err
	}

	// 4. 将统计结果转换为 map，便于快速查找
	planCounts := make(map[string]int)
	for _, result := range countResults {
		planCounts[result.PartNumberPrefix] = result.Count
	}

	// 5. 对每个生产计划应用智能分配算法
	for i := range plans {
		totalMatched := 0
		if count, exists := planCounts[plans[i].PartNumber]; exists {
			totalMatched = count
		}

		// 智能分配算法：按照 T -> T1 -> T2 -> T3 的顺序填充实际完成数量
		remainingCount := totalMatched

		// T 阶段
		plans[i].TActual = min(plans[i].TPlanned, remainingCount)
		plans[i].TUnfinished = plans[i].TPlanned - plans[i].TActual
		remainingCount -= plans[i].TActual

		// T1 阶段
		plans[i].T1Actual = min(plans[i].T1Planned, remainingCount)
		plans[i].T1Unfinished = plans[i].T1Planned - plans[i].T1Actual
		remainingCount -= plans[i].T1Actual

		// T2 阶段
		plans[i].T2Actual = min(plans[i].T2Planned, remainingCount)
		plans[i].T2Unfinished = plans[i].T2Planned - plans[i].T2Actual
		remainingCount -= plans[i].T2Actual

		// T3 阶段
		plans[i].T3Actual = min(plans[i].T3Planned, remainingCount)
		plans[i].T3Unfinished = plans[i].T3Planned - plans[i].T3Actual
		remainingCount -= plans[i].T3Actual

		// 计算总计字段
		plans[i].TotalInspected = plans[i].TActual + plans[i].T1Actual + plans[i].T2Actual + plans[i].T3Actual
		plans[i].TotalUnfinished = plans[i].TUnfinished + plans[i].T1Unfinished + plans[i].T2Unfinished + plans[i].T3Unfinished

		// 计算达成率（避免除以0）
		if plans[i].TotalPlanned > 0 {
			plans[i].AchievementRate = float64(plans[i].TotalInspected) / float64(plans[i].TotalPlanned) * 100.0
		} else {
			plans[i].AchievementRate = 0.0
		}
	}

	return plans, nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *ProductionPlanService) GetActiveProductionPlan(date time.Time, productModelID *uint, allowExceed bool) (*models.ProductionPlan, error) {
	// This method might need adjustment or removal based on new requirements,
	// but keeping it for now as it might be used elsewhere.
	// Since the model changed, this implementation is likely broken and needs to be updated if used.
	return nil, nil
}

func (s *ProductionPlanService) ImportProductionPlan(file multipart.File) ([]models.ProductionPlan, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil, err
	}
	defer f.Close()

	// Assuming the first sheet is the one we want
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		fmt.Printf("Error getting rows: %v\n", err)
		return nil, err
	}

	fmt.Printf("Total rows read: %d\n", len(rows))

	var plans []models.ProductionPlan
	var targetDate time.Time
	dateSet := false

	// Helper to parse int safely
	parseInt := func(s string) int {
		val, _ := strconv.Atoi(s)
		return val
	}

	// Skip header row
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}

		// Parse date
		planDateStr := row[4]
		var planDate time.Time
		var err error

		planDate, err = time.Parse("2006-01-02", planDateStr)
		if err != nil {
			planDate, err = time.Parse("06-01-02", planDateStr)
			if err != nil {
				continue
			}
		}

		// Validate all dates are the same
		if !dateSet {
			targetDate = planDate
			dateSet = true
		} else if !planDate.Equal(targetDate) {
			return nil, fmt.Errorf("所有记录的计划日期必须相同，发现: %s 和 %s",
				targetDate.Format("2006-01-02"), planDate.Format("2006-01-02"))
		}

		plan := models.ProductionPlan{
			MaterialCode:   row[0],
			PartNumber:     row[1],
			Type:           row[2],
			Manufacturer:   row[3],
			PlanDate:       planDate,
			ProductionLine: row[5],
		}

		// T
		if len(row) > 6 {
			plan.TPlanned = parseInt(row[6])
		}
		if len(row) > 7 {
			plan.TActual = parseInt(row[7])
		}
		plan.TUnfinished = plan.TPlanned - plan.TActual

		// T+1
		if len(row) > 9 {
			plan.T1Planned = parseInt(row[9])
		}
		if len(row) > 10 {
			plan.T1Actual = parseInt(row[10])
		}
		plan.T1Unfinished = plan.T1Planned - plan.T1Actual

		// T+2
		if len(row) > 12 {
			plan.T2Planned = parseInt(row[12])
		}
		if len(row) > 13 {
			plan.T2Actual = parseInt(row[13])
		}
		plan.T2Unfinished = plan.T2Planned - plan.T2Actual

		// T+3
		if len(row) > 15 {
			plan.T3Planned = parseInt(row[15])
		}
		if len(row) > 16 {
			plan.T3Actual = parseInt(row[16])
		}
		plan.T3Unfinished = plan.T3Planned - plan.T3Actual

		// 计算汇总统计
		plan.TotalPlanned = plan.TPlanned + plan.T1Planned + plan.T2Planned + plan.T3Planned
		plan.TotalInspected = plan.TActual + plan.T1Actual + plan.T2Actual + plan.T3Actual
		plan.TotalUnfinished = plan.TUnfinished + plan.T1Unfinished + plan.T2Unfinished + plan.T3Unfinished

		// 计算达成率（避免除以0）
		if plan.TotalPlanned > 0 {
			plan.AchievementRate = float64(plan.TotalInspected) / float64(plan.TotalPlanned) * 100
		} else {
			plan.AchievementRate = 0
		}

		// 特殊物料备注（假设在第23列，索引22）
		if len(row) > 22 {
			plan.SpecialNote = row[22]
		}

		plans = append(plans, plan)
	}

	fmt.Printf("Total plans to save: %d\n", len(plans))

	if len(plans) == 0 {
		return nil, fmt.Errorf("文件中没有找到有效的数据")
	}

	// Transaction to save
	err = s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&plans).Error
	})

	if err != nil {
		return nil, fmt.Errorf("保存失败: %v", err)
	}

	return plans, nil
}
