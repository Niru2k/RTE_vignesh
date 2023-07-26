package repository

import (
	//user defined package

	"todo/models"

	"gorm.io/gorm"
)

var Db *gorm.DB

func CreateTables() {
	Db.AutoMigrate(&models.Information{})
	Db.AutoMigrate(&models.TaskDetails{})

}

func CreateUser(user models.Information) {
	Db.Create(&user)
}

func ReadUserByEmail(user models.Information) (models.Information, error) {
	err := Db.Where("email=?", user.Email).First(&user).Error
	return user, err
}
func TaskPosting(post models.TaskDetails) error {
	err := Db.Create(&post).Error
	return err
}
func GetAllTask() ([]models.TaskDetails, error) {
	var creates []models.TaskDetails
	err := Db.Debug().Find(&creates).Error
	return creates, err
}
func GetTaskPostId(taskID string) (models.TaskDetails, error) {
	var create models.TaskDetails
	err := Db.Where("task_id=?", taskID).First(&create).Error
	return create, err
}
func UpdateTask(taskID string, updatetask models.TaskDetails) error {
	err := Db.Where("task_id=?", taskID).Save(&updatetask).Error
	return err
}
func ReadTaskPostById(taskID string) (models.TaskDetails, error) {
	var updatedtask models.TaskDetails
	err := Db.Where("task_id=?", taskID).First(&updatedtask).Error
	return updatedtask, err
}
func DeleteTaskPost(companyID string, deletejob models.TaskDetails) {
	Db.Where("task_id=?", companyID).Delete(&deletejob)
}

func GetTaskStatus(companyName string) ([]models.TaskDetails, error) {
	var create []models.TaskDetails
	err := Db.Where("status=?", companyName).Find(&create).Error
	return create, err
}
