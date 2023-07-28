package repository

import (
	//user defined package

	"todo/models"

	"gorm.io/gorm"
)
// database table creation 
func CreateTables(Db *gorm.DB) {
	Db.AutoMigrate(&models.Information{})
	Db.AutoMigrate(&models.TaskDetails{})
}
