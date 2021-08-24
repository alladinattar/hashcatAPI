package app

import (
	"database/sql"
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
	if err!=nil{
		return err
	}
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	DBSetup(db)
	//RabbitMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!=nil{
		log.Fatal("Failed connect to task queue", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err!=nil{
		log.Fatal("Failed create channel queue", err)
	}
	defer ch.Close()

	queueRepo := queue.NewQueue(ch)

	repo := repositories.NewHandshakeRepository(db)
	cracker := usecases.NewHashcat(cfg.Hashcat.Wordlist, cfg.Hashcat.Limit)
	handlerDB := handlers.NewHandshakes(repo)
	queueHandler := handlers.NewQueueHandler(repo, queueRepo)
	router := mux.NewRouter()

	router.Handle("/handshakes", handlerDB).Methods("GET")
	router.Handle("/result", handlerDB).Methods("GET")
	router.Handle("/task", queueHandler).Methods("POST")

	s := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	//Queue declare
	queueConsumer := queue.NewConsumer(repo, cracker)
	go queueConsumer.StartConsumeTasks()

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
