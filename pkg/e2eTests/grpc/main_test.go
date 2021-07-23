package e2etest_test

import (
	"context"
	"fmt"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/gateway/db/postgres"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
)

var TestEnv config.EnvVariablesProvider
var resource *dockertest.Resource
var mainDBConn *pgx.Conn
var mutexCreateDB sync.Mutex

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	resource, err = dockerPool.Run("postgres", "13.3", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"})
	if err != nil {
		panic(err)
	}

	TestEnv = config.EnvVariablesProviderImpl{}
	TestEnv = TestEnv.SetTestDBUrl(
		fmt.Sprintf(
			"postgres://postgres:postgres@localhost:%s/postgres?user=postgres&password=secret&sslmode=disable",
			resource.GetPort("5432/tcp"),
		),
	)

	var pgConn *pgx.Conn
	if err = dockerPool.Retry(func() error {
		pgConn, err = postgres.OpenConnection(TestEnv)
		if err != nil {
			return err
		}
		return pgConn.Ping(context.Background())
	}); err != nil {
		panic(err)
	}

	mainDBConn, err = postgres.OpenConnection(TestEnv)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	mainDBConn.Close(context.Background())
	err = dockerPool.Purge(resource)
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

func CreateNewTestDBAndEnv(name string) (config.EnvVariablesProvider, *pgx.Conn, *pgxpool.Pool) {

	name = strings.ToLower(name)

	mutexCreateDB.Lock()
	_, err := mainDBConn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s;", name))
	mutexCreateDB.Unlock()

	if err != nil {
		panic(err)
	}
	testDBUrl := fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%s/%s?user=postgres&password=secret&sslmode=disable",
		resource.GetPort("5432/tcp"),
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
