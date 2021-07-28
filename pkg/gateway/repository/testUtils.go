package repository

import (
	"context"
	"fmt"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/gateway/db/postgres"
	"strings"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
)

var (
	testEnv              config.EnvVariablesProvider
	dockerResource       *dockertest.Resource
	dbConn               *pgx.Conn
	dockerPool           *dockertest.Pool
	mutexCreateDB        sync.Mutex
	startPostgresMux     sync.Mutex
	startPostgresCounter int
)

func StartPostgresTestContainer() {
	startPostgresMux.Lock()
	if startPostgresCounter += 1; startPostgresCounter == 1 {
		startPostgresMux.Unlock()

		var err error
		dockerPool, err = dockertest.NewPool("")
		if err != nil {
			panic(err)
		}

		dockerResource, err = dockerPool.Run("postgres", "13.3", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"})
		if err != nil {
			panic(err)
		}

		testEnv = config.EnvVariablesProviderImpl{}
		testEnv = testEnv.SetTestDBUrl(
			fmt.Sprintf(
				"postgres://postgres:postgres@localhost:%s/postgres?user=postgres&password=secret&sslmode=disable",
				dockerResource.GetPort("5432/tcp"),
			),
		)

		var pgConn *pgx.Conn
		if err = dockerPool.Retry(func() error {
			pgConn, err = postgres.OpenConnection(testEnv)
			if err != nil {
				return err
			}
			return pgConn.Ping(context.Background())
		}); err != nil {
			panic(err)
		}

		dbConn, err = postgres.OpenConnection(testEnv)
		if err != nil {
			panic(err)
		}
	}
}

func FinishPostgresTestContainer() {
	startPostgresMux.Lock()
	if startPostgresCounter -= 1; startPostgresCounter == 0 {
		startPostgresMux.Unlock()

		dbConn.Close(context.Background())
		err := dockerPool.Purge(dockerResource)
		if err != nil {
			panic(err)
		}
	}
}

func CreateNewTestDBAndEnv(name string) (config.EnvVariablesProvider, *pgx.Conn, *pgxpool.Pool) {
	name = strings.ToLower(name)

	mutexCreateDB.Lock()
	_, err := dbConn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", name))
	mutexCreateDB.Unlock()

	if err != nil {
		panic(err)
	}
	testDBUrl := fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%s/%s?user=postgres&password=secret&sslmode=disable",
		dockerResource.GetPort("5432/tcp"),
		name,
	)

	var env config.EnvVariablesProvider = config.EnvVariablesProviderImpl{}
	env = env.SetTestDBUrl(testDBUrl)

	pgConn, err := postgres.OpenConnection(env)
	if err != nil {
		panic(err)
	}

	pgPool, _ := postgres.OpenConnectionPool(env)

	postgres.RunMigrations(env)

	return env, pgConn, pgPool
}
