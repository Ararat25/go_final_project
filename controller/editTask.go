package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Ararat25/go_final_project/model/entity"
)

// EditTask обработчик для изменения параметров задачи
func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	var newTask entity.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal(buf.Bytes(), &newTask)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	if newTask.ID == "" {
		sendErrorResponseData(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	_, err = strconv.Atoi(newTask.ID)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, "Указан не верный формат идентификатора")
		return
	}

	err = h.service.EditTask(&newTask)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK)
}

// sendSuccessResponse отправляет успешный ответ от сервера с пустым json
func sendSuccessResponse(w http.ResponseWriter, httpStatus int) {
	emptyStruct := make(map[string]string, 1)

	respBytes, err := json.Marshal(emptyStruct)
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	_, _ = w.Write(respBytes)
}
