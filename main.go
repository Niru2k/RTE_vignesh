package main

import (
	//user defined package
	dbOperations "todo/Lookup"
	"todo/driver"
	"todo/router"

	//third party package
	"github.com/gofiber/fiber/v2"
)

func main() {
	f := fiber.New()
	Db := driver.DatabaseConnection()
	dbOperations.UpdateDatabase(Db)
	router.SignupAndLogin(Db, f)
	router.UserAuthentication(Db, f)
	f.Listen(":8010")
}
