package main

import (
	"context"
	"errors"
	"local/panda-killer/cmd/config"
	"local/panda-killer/pkg/domain/usecase"
	"local/panda-killer/pkg/gateway/algorithms"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/db/postgres/sqlc"
	"local/panda-killer/pkg/gateway/repository"
	"local/panda-killer/pkg/gateway/rest"
	_ "local/panda-killer/swagger"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

//go:generate swag i -g main.go -o ../../swagger/

// @title Panda Killer >:D
// @version 1.0
// @description Um projeto para estudar diversas tecnologias, tecnicas e práticas utilizadas no desenvolvimento de WEB-APIs com uso de Go(lang).
// @description Mais informações no repo > https://github.com/RegisMelgaco/panda-killer
//
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	log.Info("Server is starting ...")

	env := config.EnvVariablesProviderImpl{}

	conn, err := waitPostgres(env)
	if err != nil {
		panic(err)
	}
	err = postgres.RunMigrations(env)
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	pgxConn, _ := postgres.OpenConnection(env)
	defer pgxConn.Close(context.Background())
	pgPool, _ := postgres.OpenConnectionPool(env)
	defer pgPool.Close()
	queries := sqlc.New(pgPool)

	accountRepo := repository.NewAccountRepo(queries)
	transferRepo := repository.NewTransferRepo(conn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.NewSessionTokenAlgorithms(env)
	router := rest.CreateRouter(
		env,
		usecase.NewAccountUsecase(accountRepo, passAlgo),
		usecase.NewTransferUsecase(transferRepo, accountRepo),
		usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo),
	)

	port, err := env.GetRestApiPort()
	if err != nil {
		panic(err)
	}

	log.Info("Server have started!")
	err = http.ListenAndServe(port, router)
	log.Fatal(err)
}

func waitPostgres(env config.EnvVariablesProvider) (*pgx.Conn, error) {
	for i := 0; i < 3; i++ {
		conn, err := postgres.OpenConnection(env)
		if err == nil {
			return conn, nil
		}
		time.Sleep(3 * time.Second)
	}
	return nil, errors.New("failed to connect to Postgres in 9 seconds")
}
