package model

import "github.com/Ararat25/go_final_project/dbManager"

type Service struct {
	DB *dbManager.SchedulerStore
}

func NewService(db *dbManager.SchedulerStore) *Service {
	return &Service{
		DB: db,
	}
}
