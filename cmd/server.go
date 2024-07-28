package main

import (
	"fmt"
	"github.com/Ararat25/go_final_project/controller"
	"github.com/Ararat25/go_final_project/tests"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	defaultPort = tests.Port
	webDir      = "../web"
	toDoPort    = "TODO_PORT"
)

func LoadEnvVars() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла, будут использованы значения по умолчанию")
	}
}

// runServer запускает сервер
func runServer() {
	port := getServerPort()
	server := getServerWithProperties(port)

	log.Printf("Сервер запущен на порту: %d", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Ошибка запуска сервера")
	}
}

// getServerWithProperties возвращает сервер с определенными свойтсвами
func getServerWithProperties(port int) *http.Server {
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", controller.NextDateHandler)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &server
}

// getServerPort возвращает порт для сервера
func getServerPort() int {
	port := defaultPort

	envPort := os.Getenv(toDoPort)
	if len(envPort) > 0 {
		eport, err := strconv.ParseInt(envPort, 10, 32)
		if err == nil {
			port = int(eport)
		}
	}

	return port
}
