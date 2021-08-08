package adapters

import (
	"bytes"
	"errors"
	"github.com/hashcatAPI/models"
	"io/ioutil"
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
	if err != nil {
		return nil, err
	}
	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		log.Println("Hashcat error")
		return []*models.Handshake{&models.Handshake{Status: "Hashcat error"}}, errors.New("Hashcat error")
	} else if status == 0 {
		if strings.Contains(out.String(), "found in potfile") {
			log.Println("found in potfile")
			result, err := ha.readPotfile(file)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		result, err := readResultFile()
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		response := models.Handshake{
			Status: "Exhausted",
		}
		return &response, nil
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
	data := strings.Split(strings.Replace(out.String(), "\n", "", 1), ":")
	response := models.Handshake{
		Password: data[3],
		SSID:     data[2],
		MAC:      data[0],
		Status:   "Already cracked",
	}
	return &response, nil
}

func readResultFile() (*models.Handshake, error) {
	file, err := os.Open("result")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return nil, err
	}

	data := strings.Split(strings.Replace(string(content), "\n", "", 1), ":")
	response := models.Handshake{
		Password: data[3],
		SSID:     data[2],
		MAC:      data[0],
		Status:   "Cracked",
	}
	return &response, nil
}

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}
