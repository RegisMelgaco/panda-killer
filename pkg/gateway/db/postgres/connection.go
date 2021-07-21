package postgres

import (
	"context"
	"local/panda-killer/cmd/config"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func OpenConnection() (*pgx.Conn, error) {
	dbUrl, err := config.GetDBUrl()
	if err != nil {
		panic(err)
	}

	return pgx.Connect(context.Background(), dbUrl)
}

func OpenConnectionPool() (*pgxpool.Pool, error) {
	dbUrl, err := config.GetDBUrl()
	if err != nil {
		panic(err)
	}

	return pgxpool.Connect(context.Background(), dbUrl)
}
