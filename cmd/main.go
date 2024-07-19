package main

import (
	"fmt"
	"github.com/Ararat25/go_final_project/dbManager"
)

func init() {
	LoadEnvVars()
}

func main() {
	_, err := dbManager.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	runServer()
}
