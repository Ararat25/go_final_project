package main

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Ошибка загрузки .env файла")
	}
}

func main() {
	runServer()
}
