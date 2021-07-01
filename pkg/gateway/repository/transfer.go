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

func (r TransferRepoImpl) CreateTransferAndUpdateAccountsBalances(ctx context.Context, newTransfer *transfer.Transfer, originNewBalance, destinationNewBalance float64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, createTransferSql, newTransfer.AccountOrigin, newTransfer.AccountDestination, newTransfer.Amount, newTransfer.CreatedAt).Scan(&newTransfer.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateAccountBalanceSql, originNewBalance, newTransfer.AccountOrigin)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, updateAccountBalanceSql, destinationNewBalance, newTransfer.AccountDestination)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)

	return err
}
