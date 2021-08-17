package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/hashcatAPI/handlers"
	"github.com/hashcatAPI/repositories"
	"github.com/hashcatAPI/usecases"
	"log"
	"net/http"
)


func Run() error {
	cfg, err := ReadConfig()
	if err!=nil{
		return err
	}
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS handshakes (id INTEGER PRIMARY KEY, mac TEXT, ssid TEXT, " +
		"password TEXT, time TEXT, enctyption TEXT, longitude TEXT, latitude TEXT, imei TEXT)")
	statement.Exec()
	repo := repositories.NewHandshakeRepository(db)
	cracker := usecases.NewHashcat(cfg.Hashcat.Wordlist, cfg.Hashcat.Limit)
	handlerCrack := handlers.NewUploadHandler(repo, cracker)
	handlerDB := handlers.NewHandshakes(repo)
	router := mux.NewRouter()

	router.Handle("/handshakes", handlerDB).Methods("GET")
	router.Handle("/handshakes", handlerDB).Methods("POST")
	router.Handle("/crack", handlerCrack).Methods("POST")

	s := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
