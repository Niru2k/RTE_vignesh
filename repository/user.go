package repository

import (
	//user defined package

	"todo/models"

	"gorm.io/gorm"
)

func TaskPosting(Db *gorm.DB, post models.TaskDetails) error {
	err := Db.Create(&post).Error
	return err
}
func GetTaskByUser(Db *gorm.DB, userID string) ([]models.TaskDetails, error) {
	var creates []models.TaskDetails
	err := Db.Where("user_id=?", userID).Find(&creates).Error

	return creates, err
}
func GetTaskStatus(Db *gorm.DB, taskstatus string, userid string) ([]models.TaskDetails, error) {
	var create []models.TaskDetails
	err := Db.Debug().Where("status=? AND user_id=?", taskstatus, userid).Find(&create).Error
	return create, err
}

func GetTaskById(Db *gorm.DB, taskID, userID string) (models.TaskDetails, error) {
	var creates models.TaskDetails
	err := Db.Where("task_id=? AND user_id=?", taskID, userID).Find(&creates).Error
	return creates, err
}

func ReadTaskPostById(Db *gorm.DB, taskID string) (models.TaskDetails, error) {
	var updatedtask models.TaskDetails
	err := Db.Where("task_id=?", taskID).First(&updatedtask).Error
	return updatedtask, err
}

func UpdateTask(Db *gorm.DB, task models.TaskDetails) error {
	return Db.Save(&task).Error
}

func DeleteTask(Db *gorm.DB, task models.TaskDetails) error {
	return Db.Delete(task).Error
}
