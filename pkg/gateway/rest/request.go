package rest

type CreateAccountRequest struct {
	Balance  int    `json:"balance"`
	Name     string `json:"name"`
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}

type CreateTransferRequest struct {
	OriginAccountID, DestinationAccountID int
	Amount                                int
}
