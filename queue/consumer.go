package queue

import (
	"encoding/json"
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
	} else {
		log.Println("Cracked " + strconv.Itoa(len(handshakes)) + " handshakes in file " + task.File)
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
		log.Println("Cracked BSSID: ", handshake.MAC, "Password: ", handshake.Password, "File: ", task.File)
		handshake.Latitude = task.Latitude
		handshake.Longitude = task.Longitude
		handshake.IMEI = task.IMEI
		handshake.File = task.File
		_, err := c.repo.Save(handshake)
		if err != nil {
			log.Println("Failed save handshake", err)
			continue
		}
		if c.handshakeExists(handshake.MAC){
			log.Println("Handshake already exists. Update password")
			err = c.repo.UpdatePasswordByMAC(handshake.MAC, handshake.Password)
			if err!=nil{
				log.Println("Failed update password of exists handshake", err)
				continue
			}
		}else{
			log.Println("Cracked original handshake")
			err = c.repo.SaveOriginHandshake(handshake)
			if err!=nil{
				log.Println("Failed save original handshake: ", err)
				continue
			}
		}
	}
}

func (c *Consumer)handshakeExists(mac string)bool{
	result, err := c.repo.GetByMAC(mac)
	if err!=nil{
		log.Println("Failed get handshake by MAC address")
	}
	log.Println("Len of array by mac: ", len(result))
	log.Println("Resutl by mac: ", result)
	if result!=nil{
		return true
	}
	return false
}