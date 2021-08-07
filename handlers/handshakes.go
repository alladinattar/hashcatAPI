package handlers

import (
	"encoding/json"
	"github.com/hashcatAPI/models"
	"log"
	"net/http"
)

type HandshakesHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewHandshakes(l *log.Logger, repository models.HandshakeRepository) *HandshakesHandler {
	return &HandshakesHandler{l, repository}
}

func (h *HandshakesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Send all handshakes")
	h.getHandshakes(w, r)
	return
}

func (h *HandshakesHandler) getHandshakes(w http.ResponseWriter, r *http.Request) {
	handshakes, err := h.handshakeRepo.GetAll()
	if err != nil {
		h.l.Println("Error when get handshakes")
		w.WriteHeader(500)
		return
	}

	result, err := json.MarshalIndent(handshakes, "", "  ")
	if err != nil {
		h.l.Println("Cannot Marshall handshakes")
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(result)
}
