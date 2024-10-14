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
var Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaGVja3N1bSI6IjBmMjdlODY0ZTA3ZmYyNDk5MDg4MzQ0ZGUzNjU4MjA1N2EzZTIxOTUwY2NjNDVlMDBlYjQ2OWE2YTQyOTk2MDAiLCJNYXBDbGFpbXMiOm51bGx9.YciQkhayLLHSl2OgHFCYxxu0LOaUHcXaMiZRcQRQeJ8"
