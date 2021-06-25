package repository

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"

	"github.com/jackc/pgx/v4"
)

type AccountRepoImpl struct {
	conn *pgx.Conn
}

func NewAccountRepo(conn *pgx.Conn) account.AccountRepo {
	return AccountRepoImpl{conn}
}

func (r AccountRepoImpl) CreateAccount(account *account.Account) error {
	err := r.conn.QueryRow(
		context.Background(),
		"INSERT INTO account(name, cpf, secret, balance) values($1, $2, $3, $4) RETURNING account_id;",
		account.Name, account.CPF, account.Secret, account.Balance,
	).Scan(&account.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r AccountRepoImpl) GetAccounts() ([]*account.Account, error) {
	rows, err := r.conn.Query(
		context.Background(),
		"SELECT account_id, name, cpf, secret, balance, created_at FROM account;",
	)
	if err != nil {
		return make([]*account.Account, 0), err
	}

	defer rows.Close()

	var accounts []*account.Account
	for i := 0; rows.Next(); i++ {
		var a account.Account
		rows.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)
		accounts = append(accounts, &a)
	}

	return accounts, rows.Err()
}
