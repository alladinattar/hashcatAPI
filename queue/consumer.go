package queue

import (
	"encoding/json"
	"fmt"
	"github.com/hashcatAPI/models"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
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
		true,        // durable
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
		false,   // auto-ack
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
			err = c.bruteHandshake(&task)
			if err != nil {
				log.Println("Failed brute file ", task.File, ". Error: ", err)
				c.repo.UpdateTaskState(&models.Handshake{File: task.File, IMEI: task.IMEI, Status: "Failed"})
			}
			d.Ack(false)
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
	} else {
		log.Println("Cracked " + strconv.Itoa(len(handshakes)) + " handshakes in file " + task.File)
		err = c.repo.UpdateTaskState(&models.Handshake{File: file.Name(), IMEI: task.IMEI, Status: "Finished"})
		if err != nil {
			log.Println("Failed update status of handshake: ", err)
		}
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
		fmt.Println(handshake)
		fmt.Println(task.File)
		handshake.Latitude = task.Latitude
		handshake.Longitude = task.Longitude
		handshake.IMEI = task.IMEI
		handshake.File = task.File
		_, err := c.repo.Save(handshake)
		if err != nil {
			log.Println("Failed save handshake", err)
			continue
		}
	}
}

/*func(c *Consumer) checkHandshake(handshake *models.Handshake)(bool){
	handshakes, _ := c.repo.GetByMAC(handshake.MAC)
	if len(handshakes)!=0{
		return false
	}
	return true
}
*/
