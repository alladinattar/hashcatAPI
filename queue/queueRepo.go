package queue

import (
	"github.com/streadway/amqp"
	"log"
)

type QueueRepo struct {
	Ch *amqp.Channel
}

func NewQueue(channel *amqp.Channel) *QueueRepo {
	return &QueueRepo{channel}
}

func (queue *QueueRepo) AddTask(task []byte) error {
	q, err := queue.Ch.QueueDeclare(
		"crackTasks", // name
		true,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Println("Failed declare queue", err)
		return err
	}

	err = queue.Ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        task,
		})
	if err != nil {
		log.Println("Failed add new task", err)
		return err
	}
	return nil
}
