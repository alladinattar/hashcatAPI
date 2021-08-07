package server

import (
	"github.com/gorilla/mux"
)

func SetupRouter(router *mux.Router) {
	router.Handle("/handshakes", nil).Methods("GET")
	router.Handle("/upload", nil).Methods("POST")
}
