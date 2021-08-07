package handlers

import (
	"encoding/json"
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
	file, err := receiveFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
	}
	log.Println("File recieved: ", file.Name())
	log.Println("Run hashcat with file ", file.Name())
	result, err := h.wpaCracker.CrackWPA(file)
	if err != nil {
		log.Println("crack tool error", err)
	}
	response, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Println("Failed marshall response", err)
	}
	w.Write(response)
	defer os.Remove(file.Name())
}

func receiveFile(r *http.Request) (*os.File, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	uploadedFile, err := ioutil.TempFile("./tempHandshakes", "shake-*.hccapx")
	defer uploadedFile.Close()
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
