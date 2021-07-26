package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var pidSlice []int

type UploadHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewUpload(l *log.Logger, repository models.HandshakeRepository) *UploadHandler {
	return &UploadHandler{l, repository}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.uploadFile(w, r)
	/*h.l.Println("Send all handshakes")
	h.getHandshakes(w, r)*/
	return
}

func (h *UploadHandler) uploadFile(w http.ResponseWriter, r *http.Request) {
	if len(r.Header["Imei"]) == 0 || len(r.Header["Date"]) == 0 {
		w.WriteHeader(204)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.l.Println(err)
	}

	file, _, err := r.FormFile("myFile")
	if err != nil {
		h.l.Println(err)
		w.WriteHeader(204)
		return
	}
	defer file.Close()
	/*imei := r.Header["Imei"]
	date := r.Header["Date"]*/
	//tempFile, err := ioutil.TempFile("tempHandshakes", imei[0]+"_"+date[0]+"-*.txt")
	uploadedFile, err := os.Create("test.hccapx")
	defer uploadedFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	//defer os.Remove(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	uploadedFile.Write(fileBytes)
	h.l.Println("File uploaded")
	hashcatCMD := exec.Command("hashcat", "-m2500", "test.hccapx", "rockyou.txt", "--outfile", "date:imei.crackes", "--outfile-format", "1,2", "-l", "10000")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err = hashcatCMD.Run()
	fmt.Println(out.String())
	if err != nil {
		log.Println(err)
	}
	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1{
		fmt.Println("Hashcat error")
		w.WriteHeader(500)
	} else {
		fmt.Println("fds")
		file, err := os.Open("date:imei.crackes")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			crackedPswd := strings.Split(scanner.Text(), ":")
			h.handshakeRepo.Save(crackedPswd[0], crackedPswd[2], crackedPswd[3])
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}
