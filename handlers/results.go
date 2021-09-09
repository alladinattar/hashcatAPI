package handlers

import (
	"encoding/json"
	"github.com/hashcatAPI/models"
	"log"
	"net/http"
)

type ResultsHandler struct {
	handshakeRepo models.HandshakeRepository
}

func NewResultsHandler(repository models.HandshakeRepository) *ResultsHandler {
	return &ResultsHandler{repository}
}

func (h *ResultsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("imei") == "" {
		w.Write([]byte("No imei header"))
		return
	}


	result, err := h.handshakeRepo.GetAllCrackedHandshakesByIMEI(r.Header.Get("imei"))
	if err!=nil{
		log.Println("Failed get handshakes by file: ", err)
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err!=nil{
		log.Println("Failed marshall: ", err)
		w.WriteHeader(500)
	}
	w.Write(data)

	log.Println("Send results for device with imei", r.Header.Get("imei"))
}

