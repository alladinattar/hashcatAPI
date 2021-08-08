package app

import (
	"github.com/gorilla/mux"
	"github.com/hashcatAPI/adapters"
	"github.com/hashcatAPI/handlers"
	"github.com/hashcatAPI/repositories"
	"log"
	"net/http"
	"time"
)

func Run() error {
	repo := repositories.NewHandshakeRepository(nil)
	cracker := adapters.NewHashcat("/usr/share/wordlists/rockyou.txt", 10000)
	handlerCrack := handlers.NewUploadHandler(repo, cracker)
	router := mux.NewRouter()
	router.Handle("/handshakes", nil).Methods("GET")
	router.Handle("/upload", handlerCrack).Methods("POST")

	s := http.Server{
		Addr:         ":9000",
		Handler:      router,
		IdleTimeout:  100 * time.Second,
		WriteTimeout: 100 * time.Second,
		ReadTimeout:  100 * time.Second,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
