package repository

import (
	"context"
	"local/panda-killer/pkg/domain/entity/account"
	"local/panda-killer/pkg/domain/entity/transfer"

	"github.com/jackc/pgx/v4"
)

type TransferRepoImpl struct {
	conn *pgx.Conn
}

func NewTransferRepo(conn *pgx.Conn) transfer.TransferRepo {
	return TransferRepoImpl{conn}
}

const (
	createTransferSql                  = "INSERT INTO transfer (origin_account, destination_account, amount, created_at) VALUES ($1, $2, $3, $4) RETURNING transfer_id;"
	updateAccountBalanceSql            = "UPDATE account SET balance = $1 WHERE account_id = $2;"
	selectTransfersContaningAccountSql = `
		select itable.transfer_id, itable.amount, itable.t_created_at, itable.account_id, itable.name, itable.cpf, itable.balance, itable.a_created_at, a.account_id, a.name, a.cpf, a.balance, a.created_at from (
			select t.transfer_id, t.destination_account, t.amount, t.created_at as t_created_at, a.account_id, a.name, a.cpf, a.balance, a.created_at as a_created_at 
			from transfer t
			inner join account as a ON t.origin_account = a.account_id) as itable
			inner join account as a on itable.destination_account = a.account_id
			where itable.account_id = $1 OR a.account_id = $1;
	`
)

func (r TransferRepoImpl) CreateTransferAndUpdateAccountsBalances(ctx context.Context, newTransfer *transfer.Transfer) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		createTransferSql,
		newTransfer.OriginAccount.ID,
		newTransfer.DestinationAccount.ID,
		newTransfer.Amount,
		newTransfer.CreatedAt,
	).Scan(&newTransfer.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateAccountBalanceSql, newTransfer.OriginAccount.Balance, newTransfer.OriginAccount.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, updateAccountBalanceSql, newTransfer.DestinationAccount.Balance, newTransfer.DestinationAccount.ID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (r TransferRepoImpl) GetTransfersCantainingAccount(ctx context.Context, accountID account.AccountID) ([]transfer.Transfer, error) {
	rows, err := r.conn.Query(ctx, selectTransfersContaningAccountSql, accountID)
	if err != nil {
		return []transfer.Transfer{}, err
	}

	var transfers []transfer.Transfer
	for rows.Next() {
		t := transfer.Transfer{OriginAccount: &account.Account{}, DestinationAccount: &account.Account{}}
		// itable.transfer_id, itable.amount, itable.t_created_at, itable.account_id, itable.name, itable.cpf, itable.balance, itable.a_created_at, a.account_id, a.name, a.cpf, a.balance, a.created_at
		err = rows.Scan(
			&t.ID,
			&t.Amount,
			&t.CreatedAt,
			&t.OriginAccount.ID,
			&t.OriginAccount.Name,
			&t.OriginAccount.CPF,
			&t.OriginAccount.Balance,
			&t.OriginAccount.CreatedAt,
			&t.DestinationAccount.ID,
			&t.DestinationAccount.Name,
			&t.DestinationAccount.CPF,
			&t.DestinationAccount.Balance,
			&t.DestinationAccount.CreatedAt,
		)
		if err != nil {
			return []transfer.Transfer{}, err
		}
		transfers = append(transfers, t)
	}

	return transfers, nil
}
