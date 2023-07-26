package main

import (
	//user defined package
	"todo/driver"
	"todo/repository"
	"todo/router"

)

func main() {
	driver.DatabaseConnection()
	repository.CreateTables()
	router.Router()
}
