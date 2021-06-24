package main

import (
	"context"
	"local/panda-killer/pkg/gateway"
	"local/panda-killer/pkg/gateway/db/postgres"
	"local/panda-killer/pkg/gateway/rest"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Server is starting ...")

	if err := waitPostgres(); err != nil {
		panic(err)
	}
	postgres.RunMigrations()

	log.Info("Server is started!")

	err := http.ListenAndServe(":8000",
		rest.CreateRouter(
			gateway.NewServer(),
		),
	)
	log.Fatal(err)
}

func waitPostgres() (err error) {
	for i := 0; i < 3; i++ {
		conn, err := postgres.OpenConnection()
		if err == nil {
			conn.Close(context.Background())
			return nil
		}
		time.Sleep(3 * time.Second)
	}
	return
}
