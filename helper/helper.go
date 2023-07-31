package helper

import (
	//inbuilt package
	"fmt"

	//third party package
	"github.com/joho/godotenv"
)

func Configure(filename string)error {
	err := godotenv.Load(filename)
	if err != nil {
		fmt.Printf("error at loading %s",filename)
		return err
	}
	return nil
}
