package handlers

import (
	"encoding/json"
	"github.com/hashcatAPI/models"
	"log"
	"net/http"
)

type ProgressHandler struct {
	handshakeRepo models.HandshakeRepository
}

func NewProgressHandler(repository models.HandshakeRepository) *ProgressHandler {
	return &ProgressHandler{repository}
}

func (h *ProgressHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Header.Get("imei") == "" {
		w.Write([]byte("No imei header"))
		return
	}

	log.Println("Send progress for device with imei", r.Header.Get("imei"))
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

func (h *ProgressHandler) getProgress(imei string) ([]*models.Handshake, error) {
	progressFiles, err := h.handshakeRepo.GetProgressByIMEI(imei)
	if err != nil {
		return nil, err
	}
	return progressFiles, nil
}


