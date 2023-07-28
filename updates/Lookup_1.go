package updates

import (
	"todo/driver"
	"todo/models"
)

func (Master) Lookup_1() {
	Db := driver.DatabaseConnection()
	Db.AutoMigrate(&models.Information{})
	Db.AutoMigrate(&models.TaskDetails{})
}
