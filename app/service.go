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

func (svc *Service) generateUUID() uuid.UUID {
	return svc.idGenerator()
}

func (svc *Service) utcTimestamp() time.Time {
	return svc.now().UTC().Truncate(time.Microsecond)
}
