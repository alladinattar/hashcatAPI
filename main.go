package main

import (
	"github.com/hashcatAPI/app"
)

func main() {
	/*db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS handshakes (id INTEGER PRIMARY KEY, mac TEXT, ssid TEXT, password TEXT)")
	statement.Exec()*/

	app.Run()
}
