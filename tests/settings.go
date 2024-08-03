package tests

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла, будут использованы значения по умолчанию")
	}

	port, err := strconv.Atoi(os.Getenv("TODO_PORT"))
	if err == nil && port != 0 {
		Port = port
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile != "" {
		DBFile = dbFile
	}
}

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = ``
