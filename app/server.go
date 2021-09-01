package app

import (
	"database/sql"
	"fmt"
	"github.com/hashcatAPI/queue"
	"github.com/hashcatAPI/usecases"
	"github.com/streadway/amqp"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashcatAPI/handlers"
	"github.com/hashcatAPI/repositories"
)

func Run() error {
	cfg, err := ReadConfig()
	if err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	DBSetup(db)
	//RabbitMQ connection
	fmt.Println("amqp://" + cfg.Queue.Login + ":" + cfg.Queue.Password + "@" + cfg.Queue.Addr + "/")
	conn, err := amqp.Dial("amqp://" + cfg.Queue.Login + ":" + cfg.Queue.Password + "@" + cfg.Queue.Addr + "/")
	if err != nil {
		log.Fatal("Failed connect to task queue:", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed create channel queue", err)
	}
	defer ch.Close()

	queueRepo := queue.NewQueue(ch)

	repo := repositories.NewHandshakeRepository(db)
	cracker := usecases.NewHashcat(cfg.Hashcat.Wordlist, cfg.Hashcat.Limit)
	handlerDB := handlers.NewProgressHandler(repo)
	queueHandler := handlers.NewQueueHandler(repo, queueRepo)
	router := mux.NewRouter()
	resultsHandler := handlers.NewResultsHandler(repo)
	router.Handle("/progress", handlerDB).Methods("GET")
	router.Handle("/results", resultsHandler).Methods("GET")
	router.Handle("/task", queueHandler).Methods("POST")

	s := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	//Queue declare
	queueConsumer := queue.NewConsumer(repo, cracker)
	for k := 0;k<cfg.queue.workers;k++{
		go queueConsumer.StartConsumeTasks(cfg.Queue.Login, cfg.Queue.Password, cfg.Queue.Addr)
	}

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
