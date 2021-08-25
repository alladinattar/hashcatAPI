package usecases

import (
	"bytes"
	"errors"
	"github.com/hashcatAPI/models"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type HashcatAdapter struct {
	wordList string
	limit    int
}

func NewHashcat(wordlist string, limit int) *HashcatAdapter {
	return &HashcatAdapter{wordlist, limit}
}

func (ha *HashcatAdapter) CrackWPA(file *os.File) (crackedHandshakes []*models.Handshake, err error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", file.Name(), ha.wordList, "-l", strconv.Itoa(ha.limit))
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	hashcatCMD.Run()
	status := exitStatus(hashcatCMD.ProcessState)
	if status == 0 {
		crackedHandshakes, err = ha.readPotfile(file)
		if err != nil {
			return nil, err
		}
		err = os.Remove(file.Name())
		if err != nil {
			log.Println("Failed when remove received file", err)
			return nil, err
		}
	} else if status == 1 {
		crackedHandshakes, err = ha.readPotfile(file)
		if err != nil {
			return nil, err
		}
	} else {
		return []*models.Handshake{}, errors.New("Hashcat error " + strconv.Itoa(status))
	}
	return crackedHandshakes, nil
}

func (ha HashcatAdapter) readPotfile(file *os.File) (crackedHandshakes []*models.Handshake, err error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), ha.wordList, "--show")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err = hashcatCMD.Run()
	if err != nil {
		return nil, err
	}
	if out.String() == "" {
		return []*models.Handshake{}, nil
	}
	log.Println(string(out.Bytes()))
	data := strings.Split(out.String(), "\n")
	for _, line := range data {
		if line == "" {
			break
		}
		separatedData := strings.Split(line, ":")
		response := models.Handshake{
			Password:   separatedData[3],
			SSID:       separatedData[2],
			MAC:        separatedData[0],
			Encryption: "WPA/WPA2",
		}
		crackedHandshakes = append(crackedHandshakes, &response)
	}
	return crackedHandshakes, nil
}

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}
