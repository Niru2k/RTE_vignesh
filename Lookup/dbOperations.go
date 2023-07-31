package dbOperations

import (
	//user defined package
	"todo/models"
	"todo/updates"

	//Inbuild package
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	//Third party package
	"gorm.io/gorm"
)

func UpdateDatabase(Db *gorm.DB) {
	var update []models.DBUpdate
	var check bool
	files, err := os.ReadDir("./updates")
	if err != nil {
		log.Println("Error :", err)
		return
	}
	Db.AutoMigrate(&models.DBUpdate{})
	Db.Find(&update)
	for i := 0; i < len(files); i++ {
		for _, value := range update {
			file := fmt.Sprintf("%s.go", value.File)
			if files[i].Name() == file {
				check = true
				break
			}
		}
		if !check && files[i].Name() != "master.go" {
			var update1 models.DBUpdate
			extension := path.Ext(files[i].Name())
			index := strings.LastIndex(files[i].Name(), extension)
			update1.File = files[i].Name()[:index]
			Db.Create(&update1)
			log.Println(update1.File, "is updated...")
			update := updates.Master{}
			update.Trigger(update1.File)
		}
		check = false
	}

}
