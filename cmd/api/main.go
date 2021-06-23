package main

import (
	"local/panda-killer/pkg/gateway/rest"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Server is starting ...")
	log.Fatal(
		http.ListenAndServe(":8000", rest.CreateRouter()),
	)
}
