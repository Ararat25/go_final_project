package controller

import (
	"github.com/Ararat25/go_final_project/model"
	"net/http"
	"time"
)

// DoneTask обработчик для отметки о выполнении задачи
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat == "" {
		err = h.service.DB.DeleteTask(idInt)
		if err != nil {
			sendErrorResponseData(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		nextDate, err := model.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
			return
		}

		task.Date = nextDate

		err = model.EditTask(task, h.service.DB)
		if err != nil {
			sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	sendSuccessResponse(w, http.StatusOK)
}
