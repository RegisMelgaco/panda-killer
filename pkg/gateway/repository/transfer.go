package repository

import (
	"context"
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
	createTransferSql       = "INSERT INTO transaction (origin_account, destination_account, amount, created_at) VALUES ($1, $2, $3, $4) RETURNING transaction_id;"
	updateAccountBalanceSql = "UPDATE account SET balance = $1 WHERE account_id = $2;"
)

func (r TransferRepoImpl) CreateTransferAndUpdateAccountsBalances(ctx context.Context, newTransfer *transfer.Transfer) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	err = tx.QueryRow(
		ctx,
		createTransferSql,
		newTransfer.OriginAccount.ID,
		newTransfer.DestinationAccount.ID,
		newTransfer.Amount,
		newTransfer.CreatedAt,
	).Scan(&newTransfer.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, updateAccountBalanceSql, newTransfer.OriginAccount.Balance, newTransfer.OriginAccount.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	_, err = tx.Exec(ctx, updateAccountBalanceSql, newTransfer.DestinationAccount.Balance, newTransfer.DestinationAccount.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}
