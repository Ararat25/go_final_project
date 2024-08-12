package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ararat25/go_final_project/errors"
	"github.com/Ararat25/go_final_project/model/entity"
)

// GetTask обработчик для получения задачи из бд
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	idString := r.FormValue("id")

	id, err := checkID(idString)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	newTask, err := h.service.GetTaskById(id)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccessResponseTask(w, http.StatusOK, newTask)
}

// sendSuccessResponseTask отправляет успешный ответ с сервера от обработчика GetTask
func sendSuccessResponseTask(w http.ResponseWriter, httpStatus int, task *entity.Task) {
	respBytes, err := json.Marshal(task)
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

// checkID проверяет валидность ID
func checkID(id string) (int, error) {
	if id == "" {
		return -1, errors.ErrIdNotSpecified
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return -1, errors.ErrInvalidIdFormat
	}

	return idInt, nil
}
