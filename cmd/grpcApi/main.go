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
	"local/panda-killer/pkg/gateway/rpc"
	"local/panda-killer/pkg/gateway/rpc/gen"
	_ "local/panda-killer/swagger"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	log.Info("Server is starting ...")

	env := config.EnvVariablesProviderImpl{}

	pgConn, err := waitPostgres(env)
	if err != nil {
		panic(err)
	}
	defer pgConn.Close(context.Background())
	postgres.RunMigrations(env)

	pgPool, _ := postgres.OpenConnectionPool(env)
	defer pgPool.Close()
	queries := sqlc.New(pgPool)

	accountRepo := repository.NewAccountRepo(queries)
	transferRepo := repository.NewTransferRepo(pgConn)
	passAlgo := algorithms.PasswordHashingAlgorithmsImpl{}
	sessionAlgo := algorithms.NewSessionTokenAlgorithms(env)
	accountUsecase := usecase.NewAccountUsecase(accountRepo, passAlgo)
	authUsecase := usecase.NewAuthUsecase(accountRepo, sessionAlgo, passAlgo)
	transferUsecase := usecase.NewTransferUsecase(transferRepo, accountRepo)

	grpcPort, err := env.GetGRPCApiPort()
	if err != nil {
		panic(err)
	}
	restPort, err := env.GetRestApiPort()
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("GRPC server failed to listen: %v", err)
	}

	s := grpc.NewServer()
	api := rpc.NewApi(accountUsecase, authUsecase, transferUsecase)
	gen.RegisterPandaKillerServer(s, api)

	go func() {
		log.Printf("GRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("GRPC server failed to serve: %v", err)
		}
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0"+grpcPort,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = gen.RegisterPandaKillerHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    restPort,
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on " + restPort)
	log.Fatalln(gwServer.ListenAndServe())
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
