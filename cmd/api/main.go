package main

import (
	"local/panda-killer/pkg/gateway/rest"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Server is starting ...")
	rest.StartServer(":9000")
}
