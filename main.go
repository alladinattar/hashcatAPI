package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/hashcatAPI/server"
	"log"
)

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS handshakes (id INTEGER PRIMARY KEY, mac TEXT, ssid TEXT, password TEXT)")
	statement.Exec()

	mux := mux.NewRouter()
	server.SetupRouter(mux)
	server.Run(mux)
}
