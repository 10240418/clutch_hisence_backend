package main

import (
	"github.com/dreamskynl/godi"
	"gorm.io/gorm"
)

var DB_CONN *gorm.DB

var SERVICE_CONTAINER godi.IGoDI
