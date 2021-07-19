package main

import (
	"database/sql"
	"github.com/hashcatAPI/handlers"
	"github.com/hashcatAPI/repositories"
	"log"
	"net/http"
	"os"
)

func main() {
	l := log.New(os.Stdout, "hshandler", log.LstdFlags)
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		l.Fatal(err)
	}

	repo := repositories.NewHandshakeRepository(db)
	hh := handlers.NewHandshakes(l, repo)
	mux := http.NewServeMux()

	mux.Handle("/handshakes", hh)
	s := http.Server{
		Addr:    ":9000",
		Handler: mux,
	}
	err = s.ListenAndServe()
	if err != nil {
		l.Fatal(err)
	}
}
