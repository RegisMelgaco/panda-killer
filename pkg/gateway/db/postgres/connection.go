package postgres

import (
	"context"
	"local/panda-killer/cmd/config"

	"github.com/jackc/pgx/v4"
)

func OpenConnection() (*pgx.Conn, error) {
	dbUrl, err := config.GetDBUrl()
	if err != nil {
		panic(err)
	}

	return pgx.Connect(context.Background(), dbUrl)
}
