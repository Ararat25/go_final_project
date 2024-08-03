package controller

import (
	"encoding/json"
	"github.com/Ararat25/go_final_project/customError"
	"github.com/Ararat25/go_final_project/dbManager"
	"log"
	"net/http"
	"strconv"
)

// GetTask обработчик для получения задачи из бд
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	idInt, err := checkID(id)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.service.DB.GetTaskById(idInt)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccessResponseTask(w, http.StatusOK, task)
}

// sendSuccessResponseTask отправляет успешный ответ с сервера от обработчика GetTask
func sendSuccessResponseTask(w http.ResponseWriter, httpStatus int, task *dbManager.Task) {
	respBytes, err := json.Marshal(task)
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}

// checkID проверяет валидность ID
func checkID(id string) (int, error) {
	if id == "" {
		return -1, customError.ErrIdNotSpecified
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return -1, customError.ErrInvalidIdFormat
	}

	return idInt, nil
}
