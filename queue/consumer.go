package queue

import (
	"github.com/streadway/amqp"
	"log"
)

func StartConsumeTasks(){
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!=nil{
		log.Fatal("Failed connect to queue", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err!=nil{
		log.Fatal("Failed create channel", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"crackTasks", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err!=nil{
		log.Fatal("Failed declare queue", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err!=nil{
		log.Fatal("Failed consume messages", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a file: %s", d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}



