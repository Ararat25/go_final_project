package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/go_final_project/dbManager"
	"github.com/Ararat25/go_final_project/model"
	"log"
	"net/http"
)

var timeLayout = "20060102"

type Response struct {
	Id    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// AddTask обработчик для добавления новых задач
func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
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

	id, err := model.AddTask(&task, h.service.DB)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccessResponseData(w, http.StatusOK, id)
}

func sendSuccessResponseData(w http.ResponseWriter, httpStatus int, id int64) {
	respBytes, err := json.Marshal(Response{Id: id})
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}

func sendErrorResponseData(w http.ResponseWriter, httpStatus int, error string) {
	respBytes, err := json.Marshal(Response{Error: error})
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(error)
	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}
