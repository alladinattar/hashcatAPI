package adapters

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashcatAPI/models"
	"log"
	"os"
	"os/exec"
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

func (ha *HashcatAdapter) CrackWPA(file *os.File) ([]*models.Handshake, error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), ha.wordList, "--outfile", "result", "--outfile-format", "1,2", "-l", string(ha.limit))
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err := hashcatCMD.Run()
	fmt.Println(out.String())
	if err != nil {
		return nil, err
	}
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
	data := strings.Split(out.String(), "\n")
	for _, line := range data {
		separatedData := strings.Split(line, ":")
		response := models.Handshake{
			Password: separatedData[3],
			SSID:     separatedData[2],
			MAC:      separatedData[0],
			Status:   "Cracked",
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
