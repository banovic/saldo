package app

type Repositories struct {
}

type UnitOfWork interface {
	Execute(work func(Repositories) error) error
}
