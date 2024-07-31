package controller

import (
	"github.com/Ararat25/go_final_project/model"
	"log"
	"net/http"
	"time"
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
	w.Write([]byte(nextDate))
}
