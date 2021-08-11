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

type CrackHandler struct {
	handshakeRepo models.HandshakeRepository
	wpaCracker    models.Cracker
}

func NewUploadHandler(repository models.HandshakeRepository, wpaCracker models.Cracker) *CrackHandler {
	return &CrackHandler{repository, wpaCracker}
}

func (h *CrackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.bruteHandshake(w, r)
	return
}

func (h *CrackHandler) bruteHandshake(w http.ResponseWriter, r *http.Request) {
	file, err := h.receiveFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("Failed receive file"))
	}
	defer file.Close()
	log.Println("Run hashcat with file ", file.Name())
	handshakes, err := h.wpaCracker.CrackWPA(file)
	if err != nil {
		log.Println("Crack tool error", err)
		return
	}
	if len(handshakes) == 0 {
		w.Write([]byte("No cracked handshakes"))
		log.Println("No cracked handshakes")
		return
	}
	result, err := json.MarshalIndent(handshakes, "", "  ")
	if err != nil {
		log.Println("Failed marshall response", err)
		return
	}
	log.Println(string(result))
	_, err = h.handshakeRepo.Save(handshakes)
	if err != nil {
		log.Println("Failed save handshake", err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(result)
}

func (h *CrackHandler) receiveFile(r *http.Request) (*os.File, error) {
	longitude := r.Header.Get("lon")
	latitude := r.Header.Get("lat")
	imei := r.Header.Get("imei")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileName := fmt.Sprintf("./tempHandshakes/shake-%s-%s-%s", imei, longitude, latitude)
	uploadedFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	uploadedFile.Write(fileBytes)
	return uploadedFile, nil
}
