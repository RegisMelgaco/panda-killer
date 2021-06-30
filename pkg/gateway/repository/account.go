package repository

import (
	"context"
	"fmt"
	"local/panda-killer/pkg/domain/entity/account"

	"github.com/jackc/pgx/v4"
)

type AccountRepoImpl struct {
	conn *pgx.Conn
}

func NewAccountRepo(conn *pgx.Conn) account.AccountRepo {
	return AccountRepoImpl{conn}
}

func (r AccountRepoImpl) CreateAccount(ctx context.Context, account *account.Account) error {
	err := r.conn.QueryRow(
		ctx,
		"INSERT INTO account(name, cpf, secret, balance) values($1, $2, $3, $4) RETURNING account_id, created_at;",
		account.Name, account.CPF, account.Secret, account.Balance,
	).Scan(&account.ID, &account.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r AccountRepoImpl) GetAccounts(ctx context.Context) ([]*account.Account, error) {
	rows, err := r.conn.Query(
		ctx,
		"SELECT account_id, name, cpf, secret, balance, created_at FROM account;",
	)
	if err != nil {
		return make([]*account.Account, 0), err
	}

	defer rows.Close()

	var accounts []*account.Account
	for rows.Next() {
		var a account.Account
		rows.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)
		accounts = append(accounts, &a)
	}

	return accounts, rows.Err()
}

func (r AccountRepoImpl) GetAccount(ctx context.Context, accountID int) (*account.Account, error) {
	row := r.conn.QueryRow(ctx, "SELECT account_id, name, cpf, secret, balance, created_at FROM account WHERE account_id = $1 FETCH FIRST ROW ONLY;", fmt.Sprint(accountID))

	var a account.Account
	err := row.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)

	return &a, err
}
