package model

import (
	"github.com/Ararat25/go_final_project/customError"
	"github.com/Ararat25/go_final_project/dbManager"
	"time"
)

// AddTask добавляет новую задачу в бд
func AddTask(task *dbManager.Task, db *dbManager.SchedulerStore) (int64, error) {
	err := checkTask(task)
	if err != nil {
		return -1, err
	}

	id, err := db.AddTask(*task)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// EditTask изменяет пареметры задачи
func EditTask(task *dbManager.Task, db *dbManager.SchedulerStore) error {
	err := checkTask(task)
	if err != nil {
		return err
	}

	err = db.EditTaskById(task)
	if err != nil {
		return err
	}

	return nil
}

// checkTask проверяет валидность переданных параметров задачи
func checkTask(task *dbManager.Task) error {
	if task.Title == "" {
		return customError.ErrTaskTitleNotSpecified
	}

	var date string

	if task.Date == "" {
		date = time.Now().Format(timeLayout)
	} else {
		taskDateParse, err := time.Parse(timeLayout, task.Date)
		if err != nil {
			return err
		}

		if task.Repeat != "" {
			repeat, err := ParseRepeat(task.Repeat)
			if err != nil {
				return err
			}

			if len(repeat.FirstSlice) != 0 && time.Now().Format(timeLayout) <= taskDateParse.Format(timeLayout) {
				return nil
			}

			if repeat.Period == "d" && repeat.FirstSlice[0] == 1 {
				date = time.Now().Format(timeLayout)
			} else {
				date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					return err
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

	return nil
}
