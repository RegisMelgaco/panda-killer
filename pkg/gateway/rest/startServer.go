package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func StartServer(address string) {
	router := chi.NewRouter()
	log.Fatal(
		http.ListenAndServe(address, router),
	)
}
