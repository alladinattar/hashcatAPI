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

func (ha *HashcatAdapter) CrackWPA(file *os.File) (*models.CrackResult, error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), ha.wordList, "--outfile", "result", "--outfile-format", "1,2", "-l", string(ha.limit))
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err := hashcatCMD.Run()
	if err != nil {
		return nil, err
	}
	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		log.Println("Hashcat error")
		return &models.CrackResult{Status: "Hashcat error"}, errors.New("Hashcat error")
	} else if status == 0 {
		if strings.Contains(out.String(), "found in potfile") {
			log.Println("Found in potfile!")
			result, err := readPotfile(file)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			result, err := readResultFile()
			if err != nil {
				return nil, err
			}
			return result, nil
		}

	} else {
		response := models.CrackResult{
			Status: "Exhausted",
		}
		return &response, nil
	}
}

func readPotfile(file *os.File) (*models.CrackResult, error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/rockyou.txt", "--show")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err := hashcatCMD.Run()
	if err != nil {
		return nil, err
	}
	data := strings.Split(strings.Replace(out.String(), "\n", "", 1), ":")
	response := models.CrackResult{
		Password: data[3],
		Ssid:     data[2],
		Mac:      data[0],
		Status:   "Cracked",
	}
	return &response, nil
}

func readResultFile() (*models.CrackResult, error) {
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
	response := models.CrackResult{
		Password: data[3],
		Ssid:     data[2],
		Mac:      data[0],
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
