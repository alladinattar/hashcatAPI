package app

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Run(router *mux.Router) error {
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
