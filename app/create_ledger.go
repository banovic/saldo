package app

type CreateLedgerRequest struct{}

type CreateLedgerResponse struct{}

func (*Service) CreateLedger(req CreateLedgerRequest) (CreateLedgerResponse, error) {
	return CreateLedgerResponse{}, nil
}
