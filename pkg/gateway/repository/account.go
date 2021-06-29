package repository

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
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
	for i := 0; rows.Next(); i++ {
		var a account.Account
		rows.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)
		accounts = append(accounts, &a)
	}

	return accounts, rows.Err()
}

func (r AccountRepoImpl) GetAccountBalance(ctx context.Context, accountID int) (float64, error) {
	rows, err := r.conn.Query(ctx, "SELECT balance FROM account WHERE account_id = $1 LIMIT 1;", accountID)
	if err != nil {
		return 0, err
	}
	logrus.Debug(rows.RawValues())
	if !rows.Next() {
		return 0, account.ErrAccountNotFound
	}

	var balance float64
	err = rows.Scan(&balance)

	return balance, err
}
