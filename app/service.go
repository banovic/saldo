package app

import "time"

type Service struct {
	uow UnitOfWork
	now func() time.Time
}

func NewService(uow UnitOfWork, now func() time.Time) *Service {
	return &Service{uow: uow, now: now}
}
