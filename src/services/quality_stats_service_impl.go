package services

import (
	"fmt"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"gorm.io/gorm"
)

type QualityStatsService struct {
	db *gorm.DB
}

func NewQualityStatsService(db *gorm.DB) (IQualityStatsService, error) {
	return &QualityStatsService{db: db}, nil
}

func (s *QualityStatsService) GetQualityStats(startDate, endDate time.Time) (*models.QualityStatsResponse, error) {
	// 获取合格率统计
	qualityRate, err := s.getQualityRateStats(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 获取不良类型分布
	defectTypeDistribution, err := s.getDefectTypeDistribution(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 获取供应商不良趋势
	supplierDefectTrend, err := s.getSupplierDefectTrend(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 获取各类型不良趋势
	defectTrendByType, err := s.getDefectTrendByType(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return &models.QualityStatsResponse{
		QualityRate:            *qualityRate,
		DefectTypeDistribution: defectTypeDistribution,
		SupplierDefectTrend:    supplierDefectTrend,
		DefectTrendByType:      *defectTrendByType,
	}, nil
}

// 获取合格率统计
func (s *QualityStatsService) getQualityRateStats(startDate, endDate time.Time) (*models.QualityRateStats, error) {
	var totalCount, defectCount int64

	// 统计总数
	if err := s.db.Model(&models.Product{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 统计不合格数
	if err := s.db.Model(&models.Product{}).
		Where("created_at BETWEEN ? AND ? AND has_defect = ?", startDate, endDate, true).
		Count(&defectCount).Error; err != nil {
		return nil, err
	}

	qualifiedCount := totalCount - defectCount
	var qualityRate float64
	if totalCount > 0 {
		qualityRate = float64(qualifiedCount) / float64(totalCount) * 100
	}

	return &models.QualityRateStats{
		QualifiedCount: int(qualifiedCount),
		DefectCount:    int(defectCount),
		TotalCount:     int(totalCount),
		QualityRate:    qualityRate,
	}, nil
}

// 获取不良类型分布
func (s *QualityStatsService) getDefectTypeDistribution(startDate, endDate time.Time) ([]models.DefectTypeItem, error) {
	var results []struct {
		DefectReason string
		Count        int64
	}

	// 查询不良产品的缺陷原因分布
	if err := s.db.Model(&models.Product{}).
		Select("defect_reason, COUNT(*) as count").
		Where("created_at BETWEEN ? AND ? AND has_defect = ? AND defect_reason != ''", startDate, endDate, true).
		Group("defect_reason").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// 获取总的不良数量
	var totalDefectCount int64
	for _, result := range results {
		totalDefectCount += result.Count
	}

	// 转换为响应格式
	var defectTypes []models.DefectTypeItem
	for _, result := range results {
		rate := float64(result.Count) / float64(totalDefectCount) * 100
		defectTypes = append(defectTypes, models.DefectTypeItem{
			Type:  result.DefectReason,
			Count: int(result.Count),
			Rate:  rate,
		})
	}

	return defectTypes, nil
}

// 获取供应商不良趋势
func (s *QualityStatsService) getSupplierDefectTrend(startDate, endDate time.Time) ([]models.SupplierDefectTrend, error) {
	// 首先获取所有供应商
	var suppliers []models.Supplier
	if err := s.db.Find(&suppliers).Error; err != nil {
		return nil, err
	}

	var supplierTrends []models.SupplierDefectTrend

	for _, supplier := range suppliers {
		dailyData, err := s.getSupplierDailyDefectData(supplier.ID, startDate, endDate)
		if err != nil {
			return nil, err
		}

		if len(dailyData) > 0 {
			supplierTrends = append(supplierTrends, models.SupplierDefectTrend{
				SupplierName: supplier.Name,
				DailyData:    dailyData,
			})
		}
	}

	return supplierTrends, nil
}

// 获取供应商每日不良率数据
func (s *QualityStatsService) getSupplierDailyDefectData(supplierID int64, startDate, endDate time.Time) ([]models.DailyDefectRate, error) {
	// 按天分组统计每个供应商的产品数量和不良数量
	var results []struct {
		Date        string
		TotalCount  int64
		DefectCount int64
	}

	query := `
		SELECT 
			DATE(p.created_at) as date,
			COUNT(*) as total_count,
			SUM(CASE WHEN p.has_defect = true THEN 1 ELSE 0 END) as defect_count
		FROM products p
		INNER JOIN product_models pm ON p.product_model_id = pm.id
		WHERE pm.supplier_id = ? 
			AND p.created_at BETWEEN ? AND ?
		GROUP BY DATE(p.created_at)
		ORDER BY date
	`

	if err := s.db.Raw(query, supplierID, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	var dailyData []models.DailyDefectRate
	for _, result := range results {
		var defectRate float64
		if result.TotalCount > 0 {
			defectRate = float64(result.DefectCount) / float64(result.TotalCount) * 100
		}

		dailyData = append(dailyData, models.DailyDefectRate{
			Date:        result.Date,
			DefectRate:  defectRate,
			TotalCount:  int(result.TotalCount),
			DefectCount: int(result.DefectCount),
		})
	}

	return dailyData, nil
}

// 获取各类型不良趋势
func (s *QualityStatsService) getDefectTrendByType(startDate, endDate time.Time) (*models.DefectTrendByType, error) {
	defectTypes := map[string]string{
		"terminal":   "端子变形",
		"tag":        "铭牌不良",
		"appearance": "外观不良",
		"noise":      "轴承噪音",
	}

	trends := &models.DefectTrendByType{}

	for key, defectType := range defectTypes {
		dailyData, err := s.getDefectTypeDailyData(defectType, startDate, endDate)
		if err != nil {
			return nil, err
		}

		switch key {
		case "terminal":
			trends.TerminalData = dailyData
		case "tag":
			trends.TagData = dailyData
		case "appearance":
			trends.AppearanceData = dailyData
		case "noise":
			trends.NoiseData = dailyData
		}
	}

	return trends, nil
}

// 获取特定缺陷类型的每日数据
func (s *QualityStatsService) getDefectTypeDailyData(defectType string, startDate, endDate time.Time) ([]models.DailyDefectCount, error) {
	var results []struct {
		Date  string
		Count int64
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM products 
		WHERE has_defect = true 
			AND defect_reason LIKE ?
			AND created_at BETWEEN ? AND ?
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	searchPattern := fmt.Sprintf("%%%s%%", defectType)
	if err := s.db.Raw(query, searchPattern, startDate, endDate).Scan(&results).Error; err != nil {
		return nil, err
	}

	var dailyData []models.DailyDefectCount
	for _, result := range results {
		dailyData = append(dailyData, models.DailyDefectCount{
			Date:  result.Date,
			Count: int(result.Count),
		})
	}

	return dailyData, nil
}
