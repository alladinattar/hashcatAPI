package queue

import (
	"encoding/json"
	"github.com/hashcatAPI/models"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Consumer struct {
	repo    models.HandshakeRepository
	cracker models.Cracker
}

func NewConsumer(repo models.HandshakeRepository, cracker models.Cracker) *Consumer {
	return &Consumer{repo, cracker}
}

func (c *Consumer) StartConsumeTasks(login, password, addr string) {
	conn, err := amqp.Dial("amqp://" + login + ":" + password + "@" + addr + "/")
	if err != nil {
		log.Fatal("Failed connect to queue", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed create channel", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"crackTasks", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("Failed declare queue", err)
	}
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatal(err)
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
	if err != nil {
		log.Fatal("Failed consume messages", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received task: %s", d.Body)
			task := models.Handshake{}
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Println("Failed unmarshall data from queue", err)
				continue
			}
			err = c.repo.AddTaskToDB(&task)
			if err != nil {
				log.Println("Failed save task to db:", err)
				continue
			}
			err = c.bruteHandshake(&task)
			if err != nil {
				log.Println("Failed brute file ", task.File, ". Error: ", err)

			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (c *Consumer) bruteHandshake(task *models.Handshake) error {
	log.Println("Run hashcat with file", task.File)
	file, err := os.Open(task.File)
	if err != nil {
		return err
	}

	err = c.repo.UpdateTaskState(&models.Handshake{File: file.Name(), IMEI: task.IMEI, Status: "In progress"})
	if err != nil {
		return err
	}

	handshakes, err := c.cracker.CrackWPA(file)
	if err != nil {
		return err
	}

	if len(handshakes) == 0 {
		log.Println("No cracked handshakes in ", file.Name())
		err := c.repo.UpdateTaskState(&models.Handshake{File: file.Name(), IMEI: task.IMEI, Status: "Finished"})
		if err != nil {
			return err
		}
		return nil
	}

	c.SaveHandshakes(handshakes, task)
	c.repo.UpdateTaskState(&models.Handshake{File: file.Name(), IMEI: task.IMEI, Status: "Finished"})
	if err != nil {
		log.Println("Failed save handshake", err)
		return err
	}
	return nil
}

func (c *Consumer) SaveHandshakes(handshakes []*models.Handshake, task *models.Handshake) {
	for _, handshake := range handshakes {
		handshake.Latitude = task.Latitude
		handshake.Longitude = task.Longitude
		handshake.IMEI = task.IMEI
		handshake.File = task.File
		check := c.checkHandshakeInDB(handshake)
		if check {
			_, err := c.repo.Save(handshake)
			if err!=nil{
				log.Println("Failed save handshake", err)
				continue
			}
			continue
		} else {
			log.Println("Invalid handshake. Incomplete device information received")
			continue
		}
	}
}

func (c *Consumer) checkHandshakeInDB(handshake *models.Handshake) bool {
	if handshake.SSID == "" || handshake.Password == "" || handshake.MAC == "" || handshake.IMEI == "" {
		return false
	}
	return true
}
