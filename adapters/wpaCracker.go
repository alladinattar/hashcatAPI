package adapters

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
	"time"
)

type HashcatAdapter struct {
	wordList string
	limit    int
}

func NewHashcat(wordlist string, limit int) *HashcatAdapter {
	return &HashcatAdapter{wordlist, limit}
}

func (ha *HashcatAdapter) CrackWPA(file *os.File) ([]*models.Handshake, error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", file.Name(), ha.wordList, "-l", strconv.Itoa(ha.limit))
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	hashcatCMD.Run()

	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		log.Println("Hashcat error")
		return []*models.Handshake{{Status: "Hashcat error"}}, errors.New("Hashcat error")
	} else {
		crackedShakes, err := ha.readPotfile(file)
		if err != nil {
			return nil, err
		}
		return crackedShakes, nil
	}
}

func (ha HashcatAdapter) readPotfile(file *os.File) ([]*models.Handshake, error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/rockyou.txt", "--show")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err := hashcatCMD.Run()
	if err != nil {
		return nil, err
	}
	if out.String() == "" {
		return []*models.Handshake{}, nil
	}
	crackedHandshakes := []*models.Handshake{}
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
			Status:     "Cracked",
			Encryption: "WPA/WPA2",
			Time:       strconv.Itoa(int(time.Now().Unix())),
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
