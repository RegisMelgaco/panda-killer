package rest

type AccountBalanceResponse struct {
	Balance float64 `json:"balance"`
}

type ErrorResponse struct {
	Message string
}

type CreatedAccountResponse struct {
	ID int `json:"id"`
}

type GetAccountResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	CPF  string `json:"cpf"`
}

type CreateTransferResponse struct {
	ID int `json:"id"`
}
