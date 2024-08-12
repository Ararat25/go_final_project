package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ararat25/go_final_project/model/entity"
)

// ResponseTasks структура для упешного ответа с сервера
type ResponseTasks struct {
	Tasks []entity.Task `json:"tasks"`
}

// ResponseError структура для ответа с сервера с текстом ошибки
type ResponseError struct {
	Error string `json:"error"`
}

// Find обработчик для получения задач из бд
func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := h.service.Find(search)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponseTasks(w, http.StatusOK, checkTasks(tasks))
}

// checkTasks если переданная структура равна nil, то возвращается пустая структура
func checkTasks(tasks []entity.Task) []entity.Task {
	if len(tasks) == 0 {
		return []entity.Task{}
	}

	return tasks
}

// sendSuccessResponseTasks отправляет успешный ответ с сервера от обработчика GetTasks
func sendSuccessResponseTasks(w http.ResponseWriter, httpStatus int, tasks []entity.Task) {
	respBytes, err := json.Marshal(ResponseTasks{Tasks: tasks})
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	_, err = w.Write(respBytes)
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}
}
