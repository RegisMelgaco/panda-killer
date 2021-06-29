package main

import (
	"errors"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Server is starting ...")

	conn, err := waitPostgres()
	if err != nil {
		panic(err)
	}
	postgres.RunMigrations()

	router := rest.CreateRouter(
		usecase.NewAccountUsecase(
			repository.NewAccountRepo(
				conn,
			),
		),
	)

	port, err := config.GetRestApiPort()
	if err != nil {
		panic(err)
	}

	log.Info("Server have started!")
	err = http.ListenAndServe(port, router)
	log.Fatal(err)
}

func waitPostgres() (*pgx.Conn, error) {
	for i := 0; i < 3; i++ {
		conn, err := postgres.OpenConnection()
		if err == nil {
			return conn, nil
		}
		time.Sleep(3 * time.Second)
	}
	return nil, errors.New("Failed to connect to Postgres in 9 seconds")
}
