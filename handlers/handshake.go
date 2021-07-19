package handlers

import (
	"encoding/json"
	"github.com/hashcatAPI/models"
	"log"
	"net/http"
)

type HandshakeHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewHandshakes(l *log.Logger, repository models.HandshakeRepository) *HandshakeHandler {
	return &HandshakeHandler{l, repository}
}

func (h *HandshakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println()
	switch r.Method {
	case http.MethodGet:
		h.getHandshakes(w, r)
		return
	default:
		w.Write([]byte("Invalid Method"))
	}
	//handshake, err := h.handshakeRepo.GetByID(1)
}

func (h *HandshakeHandler) getHandshakes(w http.ResponseWriter, r *http.Request) {
	handshakes, err := h.handshakeRepo.GetHandshakes()
	if err != nil {
		h.l.Println("Error when get handshakes")
	}
	result, err := json.MarshalIndent(handshakes, "", "  ")
	if err != nil {
		h.l.Println("Cannot Marshall handshakes")
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(result)
}
