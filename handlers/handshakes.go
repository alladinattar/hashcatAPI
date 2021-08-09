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
	case "POST":
		log.Println("Save handshake")
		var resp Response
		var handshakes []*models.Handshake
		json.NewDecoder(r.Body).Decode(&handshakes)
		defer r.Body.Close()
		err := h.saveHandshake(handshakes)
		if err != nil {
			log.Println("Cannot save new handshake", err)
			resp.Result = "Cannot save new handshake"
			resp.Comment = err.Error()
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			return
		}
		resp.Result = "Success added all handshakes"
		w.WriteHeader(201)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			w.WriteHeader(400)
			return
		}

	}

}

func (h *HandshakesHandler) getHandshakes() ([]*models.Handshake, error) {
	handshakes, err := h.handshakeRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return handshakes, nil
}

func (h *HandshakesHandler) saveHandshake(handshakes []*models.Handshake) error {
	_, err := h.handshakeRepo.Save(handshakes)
	if err != nil {
		return err
	}
	return nil
}
