package models

import "gorm.io/gorm"

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Supplier{})
	db.AutoMigrate(&ProductModel{})
	db.AutoMigrate(&ProductLine{})
	db.AutoMigrate(&Pallet{})
	db.AutoMigrate(&ProductionPlan{})
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&API{})
}
