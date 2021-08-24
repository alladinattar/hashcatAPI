package app

import (
	"database/sql"
	"fmt"
	"log"
)

func DBSetup(db *sql.DB) {
	tableWithHandshakes := fmt.Sprint("CREATE TABLE IF NOT EXISTS handshakes (id INTEGER PRIMARY KEY, mac TEXT, ssid TEXT, ",
		"password TEXT, time TEXT, enctyption TEXT, longitude TEXT, latitude TEXT, imei TEXT, file TEXT)")
	statement, _ := db.Prepare(tableWithHandshakes)
	_, err := statement.Exec()
	if err!=nil{
		log.Fatal("Failed create table with handshakes: ", err)
	}
	tableWithTasks := fmt.Sprint("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY,",
		" filename TEXT, imei TEXT, status TEXT)")
	statement, _ = db.Prepare(tableWithTasks)
	_, err = statement.Exec()
	if err!=nil{
		log.Fatal("Failed create table with tasks: ", err)
	}
}

//func QueueSetup() *amqp.Channel{
//
//
//}


