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
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/rockyou.txt", "--outfile", "result", "--outfile-format", "1,2", "-l", "10000")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err = hashcatCMD.Run()

	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		h.l.Println("Hashcat error")
		w.Write([]byte("Hashcat error"))
		w.WriteHeader(500)
	} else if status == 0 {
		if strings.Contains(out.String(), "found in potfile") {
			h.l.Println("Found in potfile")
			hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/rockyou.txt", "--show")
			var out bytes.Buffer
			hashcatCMD.Stdout = &out
			err = hashcatCMD.Run()
			if err != nil {
				h.l.Println("Failed read potfile:", err)
			}

			data := strings.Split(strings.Replace(out.String(), "\n", "", 1), ":")
			var response struct {
				Ssid     string `json:"ssid"`
				Password string `json:"password"`
				Mac      string `json:"mac"`
			}
			response.Password = data[3]
			response.Ssid = data[2]
			response.Mac = data[0]
			json.NewEncoder(w).Encode(response)
			return
		}
		file, err := os.Open("result")
		if err != nil {
			h.l.Println("No cracked handshakes")
			w.WriteHeader(200)
			return
		}
		defer os.Remove(file.Name())
		defer file.Close()

		content, _ := ioutil.ReadFile(file.Name())
		separateContent := strings.Split(string(content), ":")
		var response struct {
			Ssid     string `json:"ssid"`
			Password string `json:"password"`
			Mac      string `json:"mac"`
			Status   string `json:"status"`
		}
		response.Password = separateContent[3]
		response.Ssid = separateContent[2]
		response.Mac = separateContent[0]
		response.Status = "Cracked"
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			h.l.Println("Failed encode response:", err)
		}
	} else {
		var response struct {
			Ssid     string `json:"ssid",omitempty`
			Password string `json:"password",omitempty`
			Mac      string `json:"mac",omitempty`
			Status   string `json:"status"`
		}
		response.Status = "Exhausted"
		w.WriteHeader(200)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			h.l.Println("Failed encode response", err)
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
