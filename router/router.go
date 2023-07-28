package router

import (
	//user defined package
	"todo/authentication"
	"todo/handler"

	//third party package
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// common routes
func SignupAndLogin(Db *gorm.DB, f *fiber.App) {
	handler := handler.Database{Database: Db}
	common := f.Group("/")
	common.Post("/signup", handler.Signup)
	common.Post("/login", handler.Login)

}

//user authorization routers
func UserAuthentication(Db *gorm.DB, f *fiber.App) {
	handler := handler.Database{Database: Db}
	userauthenticated := f.Group("user/")
	userauthenticated.Use(authentication.AuthMiddleware())
	userauthenticated.Post("/posttask", handler.TaskPosting)
	userauthenticated.Put("/updatetask/:id", handler.UpdateTask)
	userauthenticated.Delete("/deletetask/:id", handler.DeleteTask)
	userauthenticated.Get("/gettasksbydate", handler.GetTasksByUserAndDate)

}
   