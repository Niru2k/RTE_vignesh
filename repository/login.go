package repository

import (
	//user defined package
	"todo/models"
	
	//third party package
	"gorm.io/gorm"
)

func CreateUser(Db *gorm.DB, user models.Information) {
	Db.Create(&user)
}

func ReadUserByEmail(Db *gorm.DB,user models.Information) (models.Information, error) {
	err := Db.Where("email=?", user.Email).First(&user).Error
	return user, err
}
