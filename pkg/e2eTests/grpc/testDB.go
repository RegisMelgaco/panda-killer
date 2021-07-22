package e2etest

import (
	"context"
	"fmt"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/gateway/db/postgres"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
)

func GetTestPgConn(ctx context.Context, env config.EnvVariablesProvider, dbName string) (*dockertest.Pool, *dockertest.Resource, config.EnvVariablesProvider, *pgx.Conn, *pgxpool.Pool, error) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, env, nil, nil, err
	}

	resource, err := dockerPool.Run("postgres", "13.3", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + dbName})
	if err != nil {
		return nil, nil, env, nil, nil, err
	}

	env = env.SetTestDBUrl(
		fmt.Sprintf(
			"postgres://postgres:postgres@localhost:%s/%s?user=postgres&password=secret&sslmode=disable",
			resource.GetPort("5432/tcp"),
			dbName,
		),
	)

	var pgConn *pgx.Conn
	if err = dockerPool.Retry(func() error {
		pgConn, err = postgres.OpenConnection(env)
		if err != nil {
			return err
		}
		return pgConn.Ping(ctx)
	}); err != nil {
		return nil, nil, env, nil, nil, err
	}

	postgres.RunMigrations(env)

	pgPool, err := postgres.OpenConnectionPool(env)
	if err != nil {
		return nil, nil, env, nil, nil, err
	}

	return dockerPool, resource, env, pgConn, pgPool, nil
}

func EraseDBArtifacts(ctx context.Context, pgPool *pgxpool.Pool, pgConn *pgx.Conn, dockerPool *dockertest.Pool, resource *dockertest.Resource) error {
	err := pgConn.Close(ctx)
	if err != nil {
		return err
	}

	pgPool.Close()

	err = dockerPool.Purge(resource)
	if err != nil {
		return err
	}

	return nil
}
