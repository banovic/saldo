package app

import (
	"time"
	"uuid"
)

type Service struct {
	uow         UnitOfWork
	now         func() time.Time
	idGenerator func() uuid.UUID
}

func NewService(uow UnitOfWork, now func() time.Time, idGenerator func() uuid.UUID) *Service {
	return &Service{uow: uow, now: now, idGenerator: idGenerator}
}
