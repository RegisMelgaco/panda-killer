package account

type AccountRepo interface {
	CreateAccount(*Account) error
	GetAccounts() ([]*Account, error)
}
