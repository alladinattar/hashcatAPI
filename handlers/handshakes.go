package handlers

import (
	"encoding/json"
	"github.com/hashcatAPI/models"
	"log"
	"net/http"
)

type HandshakesHandler struct {
	handshakeRepo models.HandshakeRepository
}

func NewHandshakes(repository models.HandshakeRepository) *HandshakesHandler {
	return &HandshakesHandler{repository}
}

func (h *HandshakesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		if r.Header.Get("imei")==""{
			w.Write([]byte("No imei header"))
			return
		}

		log.Println("Send handshakes for device with imei", r.Header.Get("imei"))
		progress, err := h.getProgress(r.Header.Get("imei"))

		if err != nil {
			log.Println("Failed get handshakes from db", err)
			w.Write([]byte("Failed get handshakes from db"))
			return
		}
		progressData, err := json.MarshalIndent(progress, "", "  ")
		if err != nil {
			log.Println("Failed marshall result", err)
		}
		w.Write(progressData)
		return
	}
}

func( h *HandshakesHandler) getProgress(imei string)([]*models.Handshake, error){
	progressFiles, err := h.handshakeRepo.GetProgressByIMEI(imei)
	if err!=nil{
		return nil, err
	}
	return progressFiles, nil
}

func (h *HandshakesHandler) getHandshakes(imei string) (string, error) {
	return "", nil
}
