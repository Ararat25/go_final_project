package controller

import (
	"encoding/json"
	"github.com/Ararat25/go_final_project/dbManager"
	"log"
	"net/http"
	"strconv"
)

// GetTask обработчик для получения задачи из бд
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		sendErrorResponseData(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, "Указан не верный формат идентификатора")
		return
	}

	task, err := h.service.DB.GetTasksById(idInt)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponseTask(w, http.StatusOK, task)
}

func sendSuccessResponseTask(w http.ResponseWriter, httpStatus int, task *dbManager.Task) {
	respBytes, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		sendErrorResponseTasks(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}
