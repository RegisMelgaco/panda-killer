package rest

type CreateAccountRequest struct {
	Balance  int    `json:"balance"`
	Name     string `json:"name"`
	CPF      string `json:"cpf" minLength:"11" maxLength:"11"`
	Password string `json:"password"`
}

type CreateTransferRequest struct {
	OriginAccountID      int `json:"origin_account_id"`
	DestinationAccountID int `json:"destination_account_id"`
	Amount               int `json:"amount"`
}

type LoginRequest struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}
