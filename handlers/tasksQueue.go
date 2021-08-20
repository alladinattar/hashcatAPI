package handlers

import (
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

func (h *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileName, err := h.receiveFile(r)
	if err != nil {
		log.Println("Failed recieve file:", err)
		fmt.Fprint(w, "Failed recieve file: ", err)
		return
	}

	err = h.Queue.AddTask(fileName)
	if err != nil {
		log.Println("Failed add task to queue", err)
		fmt.Fprint(w, "Failed add task to queue", err)
		return
	}
	log.Println("Add new task: ", fileName)

}

func (h *QueueHandler) receiveFile(r *http.Request) (string, error) {
	//fileName := r.Header.Get("filename")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	uploadedFile, err := os.Create("test")
	if err != nil {
		return "", err
	}
	uploadedFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	uploadedFile.Write(fileBytes)
	return uploadedFile.Name(), nil
}
