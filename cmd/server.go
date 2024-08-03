package main

import (
	"fmt"
	"github.com/Ararat25/go_final_project/controller"
	"github.com/Ararat25/go_final_project/dbManager"
	"github.com/Ararat25/go_final_project/middleware"
	"github.com/Ararat25/go_final_project/model"
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

// LoadEnvVars загружает переменные окружения из файла .env
func LoadEnvVars() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла, будут использованы значения по умолчанию")
	}
}

// runServer запускает сервер
func runServer() {
	db, err := dbManager.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()

	port := getServerPort()
	server := getServerWithProperties(port, db)

	log.Printf("Сервер запущен на порту: %d", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("Ошибка запуска сервера")
	}
}

// getServerWithProperties возвращает сервер с определенными свойствами
func getServerWithProperties(port int, db *dbManager.SchedulerStore) *http.Server {
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()

	service := model.NewService(db)

	handler := controller.NewHandler(service)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", controller.NextDateHandler)
	mux.Handle("/api/task", setJsonHeader(handler.Task))
	mux.Handle("/api/tasks", setJsonHeader(handler.GetTasks))
	mux.Handle("/api/task/done", setJsonHeader(handler.DoneTask))

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

// setJsonHeader устанавливает заголовок application/json в ответ
func setJsonHeader(h http.HandlerFunc) http.Handler {
	return middleware.JsonHeader(h)
}
