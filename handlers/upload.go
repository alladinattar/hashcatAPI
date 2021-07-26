package handlers

import (
	"bufio"
	"bytes"
	"errors"
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

type UploadHandler struct {
	l             *log.Logger
	handshakeRepo models.HandshakeRepository
}

func NewUpload(l *log.Logger, repository models.HandshakeRepository) *UploadHandler {
	return &UploadHandler{l, repository}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Start upload handler")
	h.uploadFile(w, r)
	return
}

func (h *UploadHandler) uploadFile(w http.ResponseWriter, r *http.Request) {
	err := recieveFile(r)
	if err!=nil{
		h.l.Println(err)
	}
	h.l.Println("File recieved")
	h.l.Println("Run hashcat...")
	hashcatCMD := exec.Command("hashcat", "-m2500", "test.hccapx", "rockyou.txt", "--outfile", "date:imei.crackes", "--outfile-format", "1,2", "-l", "10000")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err = hashcatCMD.Run()
	if err != nil {
		log.Println(err)
	}

	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1{
		fmt.Println("Hashcat error")
		w.WriteHeader(500)
	} else {
		file, err := os.Open("date:imei.crackes")
		if err != nil {
			h.l.Println("No cracked handshakes")
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		addedHash := 0
		for scanner.Scan() {
			h.l.Println(scanner.Text())
			crackedPswd := strings.Split(scanner.Text(), ":")
			count, _ := h.handshakeRepo.Save(crackedPswd[0], crackedPswd[2], "WPA",  r.Header["Imei"][0], r.Header["Date"][0], crackedPswd[3])
			addedHash+=count
		}
		if err := scanner.Err(); err != nil {
			h.l.Println(err)
		}
		h.l.Println("Added ", addedHash, "hashes")
		err = os.Remove("date:imei.crackes")
		if err!=nil{
			h.l.Println(err)
			return
		}
		h.l.Println("Temp files removed")
	}
}


func recieveFile(r *http.Request) error{
	if len(r.Header["Imei"]) == 0 || len(r.Header["Date"]) == 0 {
		//w.WriteHeader(204)
		return errors.New("No Imei or Date")
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	file, _, err := r.FormFile("myFile")
	if err != nil {
		return err
		//w.WriteHeader(204)
	}
	defer file.Close()
	/*imei := r.Header["Imei"]
	date := r.Header["Date"]*/
	//tempFile, err := ioutil.TempFile("tempHandshakes", imei[0]+"_"+date[0]+"-*.txt")
	uploadedFile, err := os.Create("test.hccapx")
	defer uploadedFile.Close()
	if err != nil {
		return err
	}
	//defer os.Remove(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	uploadedFile.Write(fileBytes)
	return nil
}
func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}
