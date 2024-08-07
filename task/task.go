package task

import (
	"time"

	"github.com/Ararat25/go_final_project/errors"
	"github.com/Ararat25/go_final_project/model"
	"github.com/Ararat25/go_final_project/model/entity"
)

var timeLayout = model.TimeLayout

// AddTask добавляет новую задачу в бд
func (s *Service) AddTask(task *entity.Task) (int64, error) {
	err := s.checkTask(task)
	if err != nil {
		return -1, err
	}

	id, err := s.db.AddTask(*task)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// GetTaskById возвращает задачу по id
func (s *Service) GetTaskById(id int) (*entity.Task, error) {
	newTask, err := s.db.GetTaskById(id)
	if err != nil {
		return nil, err
	}

	return newTask, nil
}

// GetTasks возвращает все задачи
func (s *Service) GetTasks() ([]entity.Task, error) {
	tasks, err := s.db.GetTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksByDate возвращает задачи по дате
func (s *Service) GetTasksByDate(date string) ([]entity.Task, error) {
	tasks, err := s.db.GetTasksByDate(date)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTasksBySearchString возвращает задачи по введенной строке
func (s *Service) GetTasksBySearchString(search string) ([]entity.Task, error) {
	tasks, err := s.db.GetTasksBySearchString(search)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// DoneTask отмечает задачу как выполненную
func (s *Service) DoneTask(id int) error {
	newTask, err := s.db.GetTaskById(id)
	if err != nil {
		return err
	}

	if newTask.Repeat == "" {
		err = s.db.DeleteTask(id)
		if err != nil {
			return err
		}
	} else {
		nextDate, err := model.NextDate(time.Now(), newTask.Date, newTask.Repeat)
		if err != nil {
			return err
		}

		newTask.Date = nextDate

		err = s.EditTask(newTask)
		if err != nil {
			return err
		}
	}

	return nil
}

// EditTask изменяет пареметры задачи
func (s *Service) EditTask(task *entity.Task) error {
	err := s.checkTask(task)
	if err != nil {
		return err
	}

	err = s.db.EditTaskById(task)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTask удаляет задачу
func (s *Service) DeleteTask(id int) error {
	err := s.db.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

// checkTask проверяет валидность переданных параметров задачи
func (s *Service) checkTask(task *entity.Task) error {
	if task.Title == "" {
		return errors.ErrTaskTitleNotSpecified
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
			repeat, err := model.ParseRepeat(task.Repeat)
			if err != nil {
				return err
			}

			if len(repeat.FirstSlice) != 0 && time.Now().Format(timeLayout) <= taskDateParse.Format(timeLayout) {
				return nil
			}

			if repeat.Period == "d" && repeat.FirstSlice[0] == 1 {
				date = time.Now().Format(timeLayout)
			} else {
				date, err = model.NextDate(time.Now(), task.Date, task.Repeat)
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
