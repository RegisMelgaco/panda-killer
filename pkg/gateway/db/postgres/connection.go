package postgres

import (
	"context"
	"local/panda-killer/cmd/config"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func OpenConnection(env config.EnvVariablesProvider) (*pgx.Conn, error) {
	dbUrl, err := env.GetDBUrl()
	if err != nil {
		panic(err)
	}

	return pgx.Connect(context.Background(), dbUrl)
}

func OpenConnectionPool(env config.EnvVariablesProvider) (*pgxpool.Pool, error) {
	dbUrl, err := env.GetDBUrl()
	if err != nil {
		panic(err)
	}

	return pgxpool.Connect(context.Background(), dbUrl)
}
