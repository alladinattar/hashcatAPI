package handlers

import (
	"encoding/json"
	"github.com/hashcatAPI/adapters"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type UploadHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewUpload(l *log.Logger, repository models.HandshakeRepository) *UploadHandler {
	return &UploadHandler{l, repository}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.bruteHandshake(w, r)
	return
}

func (h *UploadHandler) bruteHandshake(w http.ResponseWriter, r *http.Request) {
	file, err := receiveFile(r)
	if err != nil {
		h.l.Println(err)
	}
	h.l.Println("File recieved: ", file.Name())
	h.l.Println("Run hashcat with file ", file.Name())
	cracker := adapters.NewHashcatAdapter("/usr/share/wordlists/rockyou.txt", h.l)
	result, err := cracker.CrackWPA(file)
	if err != nil {
		h.l.Println("crack wpa error", err)
	}
	response, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		h.l.Println("Failed marshall response", err)
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
		//w.WriteHeader(204)
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
