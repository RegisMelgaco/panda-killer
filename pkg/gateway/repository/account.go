package repository

import (
	"context"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type AccountRepoImpl struct {
	conn *pgx.Conn
}

func NewAccountRepo(conn *pgx.Conn) account.AccountRepo {
	return AccountRepoImpl{conn}
}

const (
	selectAccountByCPF = `
		SELECT account_id, name, cpf, secret, balance, created_at
			FROM account
			WHERE cpf = $1
			FETCH FIRST ROW ONLY;
	`
)

func (r AccountRepoImpl) CreateAccount(ctx context.Context, a *account.Account) error {
	err := r.conn.QueryRow(
		ctx,
		"INSERT INTO account(name, cpf, secret, balance, created_at) values($1, $2, $3, $4, $5) RETURNING account_id;",
		a.Name, a.CPF, a.Secret, a.Balance, a.CreatedAt,
	).Scan(&a.ID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return account.ErrAccountCPFShouldBeUnique
		}
	}
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
	row := r.conn.QueryRow(ctx, "SELECT account_id, name, cpf, secret, balance, created_at FROM account WHERE account_id = $1 FETCH FIRST ROW ONLY;", accountID)

	var a account.Account
	err := row.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &account.Account{}, nil
	}

	return &a, nil
}

func (r AccountRepoImpl) GetAccountByCPF(ctx context.Context, cpf string) (*account.Account, error) {
	row := r.conn.QueryRow(ctx, selectAccountByCPF, cpf)

	var a account.Account
	err := row.Scan(&a.ID, &a.Name, &a.CPF, &a.Secret, &a.Balance, &a.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &account.Account{}, nil
	}

	return &a, nil
}
