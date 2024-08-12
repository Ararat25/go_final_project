package task

import "github.com/Ararat25/go_final_project/model/entity"

// Storage интерфейс хранилища для хранения задач
type Storage interface {
	Close()
	AddTask(task entity.Task) (int64, error)
	Find(filter string) ([]entity.Task, error)
	GetTaskById(id int) (*entity.Task, error)
	EditTaskById(task *entity.Task) error
	DeleteTask(id int) error
}
