package app

import "time"

type Service struct {
	uow UnitOfWork
	now func() time.Time
}
