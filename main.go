package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/hashcatAPI/handlers"
	"github.com/hashcatAPI/repositories"
	"log"
	"net/http"
	"os"
)

func main() {
	l := log.New(os.Stdout, "hshandler", log.LstdFlags)
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		l.Fatal(err)
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS handshakes (id INTEGER PRIMARY KEY, mac TEXT, ssid TEXT, password TEXT)")
    	statement.Exec()
	repo := repositories.NewHandshakeRepository(db)
	hhs := handlers.NewHandshakes(l, repo)
	uploadHandler := handlers.NewUpload(l, repo)
	mux := mux.NewRouter()

	mux.Handle("/handshakes", hhs).Methods("GET")
	mux.Handle("/upload", uploadHandler).Methods("POST")
	s := http.Server{
		Addr:    ":9000",
		Handler: mux,
	}
	err = s.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}
