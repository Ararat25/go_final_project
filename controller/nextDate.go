package controller

import (
	"log"
	"net/http"
	"time"

	"github.com/Ararat25/go_final_project/model"
)

// NextDateHandler обработчик возвращает следующий день, в зависимости от заданного правила
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowParse, err := time.Parse(timeLayout, now)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextDate, err := model.NextDate(nowParse, date, repeat)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(nextDate))
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}
}
