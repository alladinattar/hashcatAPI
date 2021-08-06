package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var counter int

type UploadHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewUpload(l *log.Logger, repository models.HandshakeRepository) *UploadHandler {
	return &UploadHandler{l, repository}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.uploadFile(w, r)
	return
}

func (h *UploadHandler) uploadFile(w http.ResponseWriter, r *http.Request) {
	file, err := receiveFile(r)
	if err != nil {
		h.l.Println(err)
	}
	h.l.Println("File recieved: ", file.Name())
	h.l.Println("Run hashcat with file ", file.Name())
	defer os.Remove(file.Name())
	h.l.Println("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/fasttrack.txt", "--outfile", "result", "--outfile-format", "1,2", "-l", "10000")
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/fasttrack.txt", "--outfile", "result", "--outfile-format", "1,2", "-l", "10000")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err = hashcatCMD.Run()
	if err != nil {
		h.l.Println("Failed run hashcat:", err)
	}
	h.l.Println(out.String())

	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		h.l.Println("Hashcat error")
		w.WriteHeader(500)
	} else {
		file, err := os.Open("result")
		defer os.Remove(file.Name())
		defer file.Close()
		if err != nil {
			h.l.Println("No cracked handshakes")
			w.WriteHeader(200)
			return
		}
		content, _ := ioutil.ReadFile(file.Name())
		separateContent := strings.Split(string(content), ":")
		var response struct {
			Ssid     string `json:"ssid"`
			Password string `json:"password"`
			Mac      string `json:"mac"`
		}
		response.Password = separateContent[3]
		response.Ssid = separateContent[2]
		response.Mac = separateContent[0]
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			h.l.Println("Failed encode response:", err)
		}
	}

}

func receiveFile(r *http.Request) (*os.File, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		return nil, err
		//w.WriteHeader(204)
	}
	defer file.Close()
	uploadedFile, err := ioutil.TempFile("./tempHandshakes", "shake-*.hccapx")
	defer uploadedFile.Close()
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

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}
