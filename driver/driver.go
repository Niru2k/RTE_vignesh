package driver

import (
	//user defined package
	"todo/helper"
	
	//built in package
	"fmt"
	"os"

	//third party package
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//connection with postgres database
func DatabaseConnection() *gorm.DB {
	err := helper.Configure(".env")
	if err != nil {
		fmt.Println("error is loading env file ")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	user := os.Getenv("USER")

	//connecting to postgres-SQL
	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	Database, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {

		fmt.Println("error in connecting with database")
	}
	fmt.Printf("%s,database connection sucessfull\n", dbname)
	return Database
}
