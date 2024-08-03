package controller

import (
	"encoding/json"
	"github.com/Ararat25/go_final_project/dbManager"
	"log"
	"net/http"
	"time"
)

var limit = 30

type ResponseTasks struct {
	Tasks []dbManager.Task `json:"tasks"`
}

type ResponseError struct {
	Error string `json:"error"`
}

// GetTasks обработчик для получения задач из бд
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	if search == "" {
		tasks, err := h.service.DB.GetTasks(limit)
		if err != nil {
			sendErrorResponseTasks(w, http.StatusInternalServerError, err.Error())
			return
		}

		sendSuccessResponseTasks(w, http.StatusOK, checkTasks(tasks))
		return
	} else {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			tasks, err := h.service.DB.GetTasksByDate(limit, date.Format(timeLayout))
			if err != nil {
				sendErrorResponseTasks(w, http.StatusInternalServerError, err.Error())
				return
			}

			sendSuccessResponseTasks(w, http.StatusOK, checkTasks(tasks))
		} else {
			tasks, err := h.service.DB.GetTasksBySearchString(limit, search)
			if err != nil {
				sendErrorResponseTasks(w, http.StatusInternalServerError, err.Error())
				return
			}

			sendSuccessResponseTasks(w, http.StatusOK, checkTasks(tasks))
		}
	}
}

func checkTasks(tasks []dbManager.Task) []dbManager.Task {
	if len(tasks) == 0 {
		return []dbManager.Task{}
	}

	return tasks
}

func sendSuccessResponseTasks(w http.ResponseWriter, httpStatus int, tasks []dbManager.Task) {
	respBytes, err := json.Marshal(ResponseTasks{Tasks: tasks})
	if err != nil {
		log.Println(err)
		sendErrorResponseTasks(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}

func sendErrorResponseTasks(w http.ResponseWriter, httpStatus int, error string) {
	respBytes, err := json.Marshal(ResponseError{Error: error})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(error)
	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}
