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

	result, err := h.getResults(r.Header.Get("imei"))
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

func (h *ResultsHandler) getResults(imei string) (map[string][]models.Handshake, error) {
	result := make(map[string][]models.Handshake)
	allHandshakes, err := h.handshakeRepo.GetAll()
	if err!=nil{
		return nil, err
	}
	for _, handshake := range allHandshakes{
		if handshake.IMEI != imei{
			continue
		}
		var handshakeTmp = models.Handshake{MAC: handshake.MAC, SSID: handshake.SSID, Password: handshake.Password, Latitude: handshake.Latitude, Longitude: handshake.Longitude}
		result[handshake.File] = append(result[handshake.File], handshakeTmp)
	}
	return result, nil
}
