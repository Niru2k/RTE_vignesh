package main

import (
	//user defined package
	"todo/driver"
	"todo/repository"
	"todo/router"
	"github.com/gofiber/fiber/v2"

)

func main() {
	f := fiber.New()
	Db:=driver.DatabaseConnection()
	repository.CreateTables(Db)
	router.SignupAndLogin(Db,f)
	router.UserAuthentication(Db,f)
	f.Listen(":3000")
}

