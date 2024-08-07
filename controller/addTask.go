package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ararat25/go_final_project/model"
	"github.com/Ararat25/go_final_project/model/entity"
)

var timeLayout = model.TimeLayout

// Response структура для ответа от сервера
type Response struct {
	Id    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// AddTask обработчик для добавления новых задач
func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
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

	id, err := h.service.AddTask(&newTask)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponseData(w, http.StatusOK, id)
}

// sendSuccessResponseData отправляет успешный ответ с сервера от обработчика AddTask
func sendSuccessResponseData(w http.ResponseWriter, httpStatus int, id int64) {
	respBytes, err := json.Marshal(Response{Id: id})
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	_, _ = w.Write(respBytes)
}

// sendErrorResponseData отпраляет ответ с текстом ошибки с сервера
func sendErrorResponseData(w http.ResponseWriter, httpStatus int, error string) {
	respBytes, err := json.Marshal(Response{Error: error})
	if err != nil {
		log.Println(err)
		log.Println(error)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(error)
	w.WriteHeader(httpStatus)
	_, _ = w.Write(respBytes)
}
