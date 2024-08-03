package model

import "github.com/Ararat25/go_final_project/dbManager"

// Service структура для хранения ссылки на бд
type Service struct {
	DB *dbManager.SchedulerStore
}

// NewService создание нового объекта Service
func NewService(db *dbManager.SchedulerStore) *Service {
	return &Service{
		DB: db,
	}
}
