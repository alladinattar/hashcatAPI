package handlers

import (
	"fmt"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
)

type UploadHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewUpload(l *log.Logger, repository models.HandshakeRepository) *UploadHandler {
	return &UploadHandler{l, repository}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("myFile")
	if err != nil {
		h.l.Println(err)
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("tempHandshakes", "upload.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	/*h.l.Println("Send all handshakes")
	h.getHandshakes(w, r)*/
	return
}

func (h *HandshakesHandler) uploadFile(w http.ResponseWriter, r *http.Request) {

}
