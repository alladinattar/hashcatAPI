package handlers

import (
	"github.com/hashcatAPI/models"
	"log"
	"net/http"
)

type HandshakeHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepo
}

func NewHandshakeHandler(l *log.Logger, repository models.HandshakeRepo) *HandshakeHandler {
	return &HandshakeHandler{l, repository}
}

func (hh *HandshakeHandler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	handshake, err := hh.handshakeRepo.GetByID(1)
}
