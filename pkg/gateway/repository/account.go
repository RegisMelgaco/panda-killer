package repository

import (
	"context"
	"errors"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/shared"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type AccountRepoImpl struct {
	q *sqlc.Queries
}

func NewAccountRepo(q *sqlc.Queries) account.AccountRepo {
	return AccountRepoImpl{q}
}

func (r AccountRepoImpl) CreateAccount(ctx context.Context, a *account.Account) error {
	accountID, err := r.q.InsertAccount(ctx, sqlc.InsertAccountParams{
		Name:      a.Name,
		Cpf:       a.CPF,
		Secret:    a.Secret,
		Balance:   int32(a.Balance),
		CreatedAt: a.CreatedAt,
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return account.ErrAccountCPFShouldBeUnique
		}
	}
	if err != nil {
		return err
	}

	a.ID = account.AccountID(accountID)

	return nil
}

func (r AccountRepoImpl) GetAccounts(ctx context.Context) ([]account.Account, error) {
	queriedAccounts, err := r.q.ListAccounts(ctx)
	if err != nil {
		return make([]account.Account, 0), err
	}

	var accounts []account.Account
	for _, a := range queriedAccounts {
		accounts = append(accounts, *toEntity(a))
	}

	return accounts, nil
}

func (r AccountRepoImpl) GetAccount(ctx context.Context, accountID account.AccountID) (*account.Account, error) {
	a, err := r.q.GetAccount(ctx, int32(accountID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, account.ErrAccountNotFound
	}

	return toEntity(a), nil
}

func (r AccountRepoImpl) GetAccountByCPF(ctx context.Context, cpf string) (*account.Account, error) {
	a, err := r.q.SelectAccountByCPF(ctx, cpf)
	if errors.Is(err, pgx.ErrNoRows) {
		return &account.Account{}, nil
	}

	return toEntity(a), nil
}

func toEntity(a sqlc.Account) *account.Account {
	return &account.Account{
		ID:        account.AccountID(a.AccountID),
		Name:      a.Name,
		CPF:       a.Cpf,
		Secret:    a.Secret,
		Balance:   shared.Money(a.Balance),
		CreatedAt: a.CreatedAt,
	}
}
