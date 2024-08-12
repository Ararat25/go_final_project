package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ararat25/go_final_project/cmd/config"
	"github.com/Ararat25/go_final_project/controller"
	"github.com/Ararat25/go_final_project/middleware"
	"github.com/Ararat25/go_final_project/task"
	"github.com/Ararat25/go_final_project/task/dbManager"
	"github.com/joho/godotenv"
)

var (
	defaultPort = 7540
	tokenTTL    = time.Hour * 8
	webDir      = "../web"
	toDoPort    = "TODO_PORT"
	tokenSalt   = "TOKEN_SALT"
	dbFile      = "TODO_DBFILE"
	password    = "TODO_PASSWORD"
)

// LoadEnvVars загружает переменные окружения из файла .env
func LoadEnvVars() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("Ошибка загрузки .env файла, будут использованы значения по умолчанию: %v\n", err)
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
func getServerWithProperties(port int, db task.Storage) *http.Server {
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()

	configData := config.NewConfig(port, []byte(os.Getenv(tokenSalt)), os.Getenv(dbFile), os.Getenv(password))

	service := task.NewService(&db, tokenTTL, configData)

	handler := controller.NewHandler(service)

	newMiddleware := middleware.NewMiddleware(service)

	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", controller.NextDateHandler)
	mux.Handle("/api/task", newMiddleware.Auth(setJsonHeader(handler.Task)))
	mux.Handle("/api/tasks", newMiddleware.Auth(setJsonHeader(handler.Find)))
	mux.Handle("/api/task/done", newMiddleware.Auth(setJsonHeader(handler.DoneTask)))
	mux.Handle("/api/signin", setJsonHeader(handler.SignIn))

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
