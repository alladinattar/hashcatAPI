package handlers

import (
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	h.l.Println("Run hashcat ...")
	defer os.Remove(file.Name())
	/*hashcatCMD := exec.Command("hashcat", "-m2500", fileName, "/usr/share/wordlists/fasttrack.txt", "--outfile", "result", "--outfile-format", "1,2", "-l", "10000")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err = hashcatCMD.Run()
	if err != nil {
		log.Println(err)
	}

	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		fmt.Println("Hashcat error")
		w.WriteHeader(500)
	} else {
		file, err := os.Open(fileName + ".crackes")
		if err != nil {
			h.l.Println("No cracked handshakes")
			w.WriteHeader(200)
			return
		}
		defer file.Close()*/

	//scanner := bufio.NewScanner(file)
	//addedHash := 0
	/*for scanner.Scan() {
		crackedPswd := strings.Split(scanner.Text(), ":")
		count, _ := h.handshakeRepo.Save(crackedPswd[0], crackedPswd[2], "WPA",  r.Header["Imei"][0], r.Header["Date"][0], crackedPswd[3])
		addedHash+=count
	}*/
	/*if err := scanner.Err(); err != nil {
		h.l.Println(err)
	}
	h.l.Println("Added ", addedHash, "hashes")
	err = os.Remove("date:imei.crackes")
	if err != nil {
		h.l.Println(err)
		return
	}
	h.l.Println("Temp files removed")*/

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
