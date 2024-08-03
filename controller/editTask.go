package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/go_final_project/dbManager"
	"github.com/Ararat25/go_final_project/model"
	"log"
	"net/http"
	"strconv"
)

// EditTask обработчик для изменения параметров задачи
func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	var task dbManager.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.ID == "" {
		sendErrorResponseData(w, http.StatusBadRequest, "Не указан идентификатор")
		return
	}

	_, err = strconv.Atoi(task.ID)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, "Указан не верный формат идентификатора")
		return
	}

	err = model.EditTask(&task, h.service.DB)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK)
}

func sendSuccessResponse(w http.ResponseWriter, httpStatus int) {
	emptyStruct := make(map[string]string, 1)

	respBytes, err := json.Marshal(emptyStruct)
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}
