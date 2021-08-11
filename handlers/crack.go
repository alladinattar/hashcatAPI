package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CrackHandler struct {
	handshakeRepo models.HandshakeRepository
	wpaCracker    models.Cracker
}

func NewUploadHandler(repository models.HandshakeRepository, wpaCracker models.Cracker) *CrackHandler {
	return &CrackHandler{repository, wpaCracker}
}

func (h *CrackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.bruteHandshake(w, r)
	return
}

func (h *CrackHandler) bruteHandshake(w http.ResponseWriter, r *http.Request) {
	file, err := h.receiveFile(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("Failed receive file"))
	}
	defer file.Close()
	log.Println("Run hashcat with file ", file.Name())
	handshakes, err := h.wpaCracker.CrackWPA(file)
	if err != nil {
		log.Println("Crack tool error", err)
		return
	}
	if len(handshakes) == 0 {
		w.Write([]byte("No cracked handshakes"))
		log.Println("No cracked handshakes")
		return
	}
	longitude := r.Header.Get("lon")
	latitude := r.Header.Get("lat")
	imei := r.Header.Get("imei")
	for _, handshake := range handshakes {
		handshake.Latitude = latitude
		handshake.Longitude = longitude
		handshake.IMEI = imei
		handshake.Time = strconv.Itoa(int(time.Now().Unix()))
	}
	result, err := json.MarshalIndent(handshakes, "", "  ")
	if err != nil {
		log.Println("Failed marshall response", err)
		return
	}
	log.Println(string(result))
	var saved int
	for _, handshake := range handshakes {
		check, err := h.CheckHandshakeInDB(handshake)
		if err != nil {
			log.Println("Failed check handshake")
			handshake.Status = "Failed check"
			continue
		}
		if check {
			_, err = h.handshakeRepo.Save(handshake)
			saved += 1
			handshake.Status = "Saved"
			continue
		} else {
			log.Println("Invalid handshake. Incomplete device information received")
			handshake.Status = "Invalid handshake. Incomplete device information received"
			continue
		}
	}
	if err != nil {
		log.Println("Failed save handshake", err)
		w.Write([]byte(err.Error()))
		return
	}
	total := fmt.Sprintf("Saved: %s\nDiscarded: %s", saved, len(handshakes)-saved)
	w.Write(result)
	w.Write([]byte(total))
}

func (h *CrackHandler) receiveFile(r *http.Request) (*os.File, error) {
	longitude := r.Header.Get("lon")
	latitude := r.Header.Get("lat")
	imei := r.Header.Get("imei")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileName := fmt.Sprintf("./tempHandshakes/shake-%s-%s-%s-%s", imei, longitude, latitude, strconv.Itoa(int(time.Now().Unix())))
	uploadedFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	uploadedFile.Write(fileBytes)
	return uploadedFile, nil
}

func (h *CrackHandler) CheckHandshakeInDB(handshake *models.Handshake) (bool, error) {
	if handshake.SSID == "" || handshake.Password == "" || handshake.MAC == "" || handshake.Latitude == "" || handshake.Longitude == "" || handshake.IMEI == "" {
		return false, nil
	} else {
		handshakes, err := h.handshakeRepo.GetByMAC(handshake.MAC)
		if err != nil {
			log.Println(err)
			return false, err
		}
		if len(handshakes) != 0 {
			log.Println("Handshake ", handshake.MAC, "found in db")
			return false, nil
		}
	}
	return true, nil
}
