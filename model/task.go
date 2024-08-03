package model

import (
	"github.com/Ararat25/go_final_project/customError"
	"github.com/Ararat25/go_final_project/dbManager"
	"time"
)

// AddTask добавляет новую задачу в бд
func AddTask(task *dbManager.Task, db *dbManager.SchedulerStore) (int64, error) {
	if task.Title == "" {
		return -1, customError.ErrTaskTitleNotSpecified
	}

	var date string

	if task.Date == "" {
		date = time.Now().Format(timeLayout)
	} else {
		taskDateParse, err := time.Parse(timeLayout, task.Date)
		if err != nil {
			return -1, err
		}

		if task.Repeat != "" {
			repeat, err := ParseRepeat(task.Repeat)
			if err != nil {
				return -1, err
			}

			if len(repeat.FirstSlice) != 0 && time.Now().Format(timeLayout) <= taskDateParse.Format(timeLayout) {
				id, err := db.AddTask(*task)
				if err != nil {
					return -1, err
				}

				return id, nil
			}

			if repeat.Period == "d" && repeat.FirstSlice[0] == 1 {
				date = time.Now().Format(timeLayout)
			} else {
				date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					return -1, err
				}
			}
		} else {
			date = taskDateParse.Format(timeLayout)

			if taskDateParse.Before(time.Now()) {
				date = time.Now().Format(timeLayout)
			}
		}
	}

	task.Date = date

	id, err := db.AddTask(*task)
	if err != nil {
		return -1, err
	}

	return id, nil
}
