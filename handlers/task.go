package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type QueueHandler struct {
	handshakeRepo models.HandshakeRepository
	Queue         models.TasksQueue
}

func NewQueueHandler(repository models.HandshakeRepository, queue models.TasksQueue) *QueueHandler {
	return &QueueHandler{repository, queue}
}

const handshakesDir = "./tempHandshakes/"
const handshakeExtension = ".hccapx"

func (h *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filename, err := h.receiveFile(r)
	if err != nil {
		log.Println("Failed recieve file:", err)
		fmt.Fprint(w, "Failed recieve file: ", err)
		return
	}
	if r.Header.Get("imei") == "" {
		fmt.Fprint(w, "Empty imei field")
		return
	}else if r.Header.Get("filename") == ""{
		fmt.Fprint(w, "Empty filename field")
		return
	}

	var task = models.Handshake{File: filename, IMEI: r.Header.Get("imei"), Status: "Queue", Latitude: r.Header.Get("latitude"),
		Longitude: r.Header.Get("longitude")}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		log.Println("Failed marshall task", err)
		w.WriteHeader(500)
	}

	err = h.Queue.AddTask(data)
	if err != nil {
		log.Println("Failed add task to queue", err)
		w.WriteHeader(500)
		return
	}
	err = h.handshakeRepo.AddTaskToDB(&task)
	if err != nil {
		log.Println("Failed save task to db:", err)
	}
	log.Println("Added new task: ", filename)
}

func (h *QueueHandler) receiveFile(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	filename := handshakesDir + r.Header.Get("filename") + handshakeExtension
	uploadedFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	fmt.Println("Recieved file:", filename)
	defer uploadedFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	uploadedFile.Write(fileBytes)
	return filename, nil
}
