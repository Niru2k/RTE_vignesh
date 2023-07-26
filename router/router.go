package router

import (
	//user defined package
	// "echo/authentication"
	"todo/authentication"
	"todo/handler"

	//third party package
	"github.com/gofiber/fiber/v2"
)

func Router() {
	f := fiber.New()

	f.Post("/signup", handler.Signup)
	f.Post("/login", handler.Login)
	f.Post("/posttask", authentication.AuthMiddleware(),handler.TaskRemainder)
	f.Get("/getalltask", authentication.AuthMiddleware(),handler.GetAllTaskDetails)
	f.Get("/getalltaskbyid/:id", authentication.AuthMiddleware(),handler.GetTaskDetailsByID)
	f.Put("/updatetaskbyid/:id", authentication.AuthMiddleware(),handler.UpdateTask)

	f.Delete("/deletetaskbyid/:id",authentication.AuthMiddleware(),handler.DeleteTask)
	f.Get("/gettaskbystatus/:status",authentication.AuthMiddleware(),handler.GetTaskStatus)
	f.Listen(":3000")
}
