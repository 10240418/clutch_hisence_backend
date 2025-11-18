package main

import (
	"fmt"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/clutchtechnology/hisense-vmi-dataserver/src/services"
	"github.com/dreamskynl/godi"
	"gorm.io/gorm"
)

func checkAdmin(db *gorm.DB) error {
	tb := db.Model(&models.User{})

	result := tb.Where("mobile = ?", "admin").Find(&[]models.User{})
	if result.RowsAffected != 0 {
		return fmt.Errorf("email already exists")
	}

	user := models.User{
		Username: "admin",
		Email:    "admin@admin.admin",
		Mobile:   "admin",
		Password: "admin",
	}

	result = tb.Create(&user)
	fmt.Println(result.RowsAffected)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func InitGodi() {
	SERVICE_CONTAINER = godi.New()

	if err := SERVICE_CONTAINER.Register(&services.APIService{}, services.NewAPIService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.PalletService{}, services.NewPalletService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.ProductLineService{}, services.NewProductLineService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.ProductModelService{}, services.NewProductModelService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.ProductService{}, services.NewProductService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.ProductionPlanService{}, services.NewProductionPlanService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.SupplierService{}, services.NewSupplierService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.UserService{}, services.NewUserService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.JwtService{}, services.NewJWTService); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.KeyManagementService{}, services.NewKeyManagementService); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.QualityStatsService{}, services.NewQualityStatsService, DB_CONN); err != nil {
		panic(err)
	}

	if err := SERVICE_CONTAINER.Register(&services.DataReportService{}, services.NewDataReportService, DB_CONN); err != nil {
		panic(err)
	}
}
