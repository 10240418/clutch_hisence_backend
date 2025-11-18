package models

// Product 对应 'Product' 表
type Product struct {
	ModelFields      `s2m:"-"`
	SN               string          `gorm:"type:char(32)" json:"sn"`
	BatchNumber      string          `gorm:"type:char(8)" json:"batchNumber,omitempty"` // 生产批次，长度为8bytes（4bytes实际+4bytes备用）
	ProductModelID   *uint           `json:"productModelId"`
	ProductModel     *ProductModel   `gorm:"foreignKey:ProductModelID" json:"productModel" s2m:"-"`
	ProductLineID    *uint           `json:"productLineId"`
	ProductLine      *ProductLine    `gorm:"foreignKey:ProductLineID" json:"productLine" s2m:"-"`
	ProductionPlanID *uint           `json:"productionPlanId"`
	ProductionPlan   *ProductionPlan `gorm:"foreignKey:ProductionPlanID" json:"productionPlan" s2m:"-"`
	PalletID         *uint           `json:"palletId"`
	Pallet           *Pallet         `gorm:"foreignKey:PalletID" json:"pallet" s2m:"-"`
	HasDefect        bool            `gorm:"default:false" json:"hasDefect"`          // 是否有缺陷
	DefectReason     string          `gorm:"type:text" json:"defectReason,omitempty"` // 缺陷原因
}

// 数据统计相关结构体
type QualityStatsQuery struct {
	StartDate string `form:"startDate" json:"startDate" binding:"required"`
	EndDate   string `form:"endDate" json:"endDate" binding:"required"`
}

type QualityStatsResponse struct {
	QualityRate            QualityRateStats      `json:"qualityRate"`
	DefectTypeDistribution []DefectTypeItem      `json:"defectTypeDistribution"`
	SupplierDefectTrend    []SupplierDefectTrend `json:"supplierDefectTrend"`
	DefectTrendByType      DefectTrendByType     `json:"defectTrendByType"`
}

type QualityRateStats struct {
	QualifiedCount int     `json:"qualifiedCount"`
	DefectCount    int     `json:"defectCount"`
	TotalCount     int     `json:"totalCount"`
	QualityRate    float64 `json:"qualityRate"`
}

type DefectTypeItem struct {
	Type  string  `json:"type"`
	Count int     `json:"count"`
	Rate  float64 `json:"rate"`
}

type SupplierDefectTrend struct {
	SupplierName string            `json:"supplierName"`
	DailyData    []DailyDefectRate `json:"dailyData"`
}

type DailyDefectRate struct {
	Date        string  `json:"date"`
	DefectRate  float64 `json:"defectRate"`
	TotalCount  int     `json:"totalCount"`
	DefectCount int     `json:"defectCount"`
}

type DefectTrendByType struct {
	TerminalData   []DailyDefectCount `json:"terminalData"`   // 电阻不良
	TagData        []DailyDefectCount `json:"tagData"`        // 反电动势
	AppearanceData []DailyDefectCount `json:"appearanceData"` // 外观
	NoiseData      []DailyDefectCount `json:"noiseData"`      // 噪音
}

type DailyDefectCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
