package rest

type CreateTransferRequest struct {
	OriginAccountID, DestinationAccountID int
	Amount                                int
}
