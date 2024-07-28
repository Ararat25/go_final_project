package main

import (
	"github.com/Ararat25/go_final_project/dbManager"
	"log"
)

func init() {
	LoadEnvVars()
}

func main() {
	db, err := dbManager.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()

	runServer()
}
