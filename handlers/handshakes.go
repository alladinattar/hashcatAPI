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
		log.Println("Send all handshakes")
		handshakes, err := h.getHandshakes()
		if err != nil {
			log.Println("Failed get handshakes from db", err)
			var resp Response
			resp.Result = "Failed get handshakes from db"
			resp.Comment = err.Error()
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(400)
			}
			return
		}
		result, err := json.MarshalIndent(handshakes, "", "  ")
		if err != nil {
			log.Println("Failed marshall result", err)
		}
		w.Write(result)
		return
	}
}

func (h *HandshakesHandler) getHandshakes() ([]*models.Handshake, error) {
	handshakes, err := h.handshakeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return handshakes, nil
}
