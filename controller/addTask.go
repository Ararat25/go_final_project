package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/go_final_project/dbManager"
	"github.com/Ararat25/go_final_project/model"
	"net/http"
	"time"
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
		setErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		setErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		setErrorResponseData(w, http.StatusBadRequest, "Не указан заголовок задачи")
		return
	}

	var date string

	if task.Date == "" {
		date = time.Now().Format(timeLayout)
	} else {
		_, err := time.Parse(timeLayout, task.Date)
		if err != nil {
			setErrorResponseData(w, http.StatusBadRequest, err.Error())
			return
		}

		if task.Repeat != "" {
			repeat, err := model.ParseRepeat(task.Repeat)
			if err != nil {
				setErrorResponseData(w, http.StatusInternalServerError, err.Error())
				return
			}

			if repeat.Period == "d" && repeat.FirstSlice[0] == 1 {
				date = time.Now().Format(timeLayout)
			} else {
				date, err = model.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					setErrorResponseData(w, http.StatusInternalServerError, err.Error())
					return
				}
			}
		} else {
			date = time.Now().Format(timeLayout)
		}
	}

	task.Date = date

	id, err := h.service.DB.AddTask(task)
	if err != nil {
		return
	}

	setSuccessResponseData(w, http.StatusOK, id)
}

func setSuccessResponseData(w http.ResponseWriter, httpStatus int, id int64) {
	respBytes, err := json.Marshal(Response{Id: id})
	if err != nil {
		setErrorResponseData(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}

func setErrorResponseData(w http.ResponseWriter, httpStatus int, error string) {
	respBytes, err := json.Marshal(Response{Error: error})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(httpStatus)
	w.Write(respBytes)
}
