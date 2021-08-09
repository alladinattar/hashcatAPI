package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/hashcatAPI/adapters"
	"github.com/hashcatAPI/handlers"
	"github.com/hashcatAPI/repositories"
	"log"
	"net/http"
	"time"
)

func Run() error {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS handshakes (id INTEGER PRIMARY KEY, mac TEXT, ssid TEXT, password TEXT, time TEXT, enctyption TEXT)")
	statement.Exec()
	repo := repositories.NewHandshakeRepository(db)
	cracker := adapters.NewHashcat("/usr/share/wordlists/rockyou.txt", 10000)
	handlerCrack := handlers.NewUploadHandler(repo, cracker)
	handlerDB := handlers.NewHandshakes(repo)
	router := mux.NewRouter()

	router.Handle("/handshakes", handlerDB).Methods("GET")
	router.Handle("/handshakes", handlerDB).Methods("POST")
	router.Handle("/upload", handlerCrack).Methods("POST")

	s := http.Server{
		Addr:         ":9000",
		Handler:      router,
		WriteTimeout: 100 * time.Second,
		ReadTimeout:  100 * time.Second,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
