package adapters

import (
	"bytes"
	"github.com/hashcatAPI/models"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type HashcatAdapter struct {
	l        *log.Logger
	wordList string
}

func NewHashcatAdapter(wordlist string, logger *log.Logger) *HashcatAdapter {
	return &HashcatAdapter{logger, wordlist}
}

func (ha *HashcatAdapter) CrackWPA(file *os.File) (*models.CrackResult, error) {
	hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), ha.wordList, "--outfile", "result", "--outfile-format", "1,2", "-l", "10000")
	var out bytes.Buffer
	hashcatCMD.Stdout = &out
	err := hashcatCMD.Run()

	if status := exitStatus(hashcatCMD.ProcessState); status != 0 && status != 1 {
		ha.l.Println("Hashcat error")
		return &models.CrackResult{Status: "Hashcat error"}, nil

	} else if status == 0 {
		if strings.Contains(out.String(), "found in potfile") {
			ha.l.Println("Found in potfile")
			hashcatCMD := exec.Command("hashcat", "-m2500", "./"+file.Name(), "/usr/share/wordlists/rockyou.txt", "--show")
			var out bytes.Buffer
			hashcatCMD.Stdout = &out
			err = hashcatCMD.Run()
			if err != nil {
				ha.l.Println("Failed read potfile:", err)
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
		file, err := os.Open("result")
		if err != nil {
			ha.l.Println("Failed open result", err)
			return nil, err
		}
		defer os.Remove(file.Name())
		defer file.Close()

		content, _ := ioutil.ReadFile(file.Name())
		data := strings.Split(string(content), ":")
		response := models.CrackResult{
			Password: data[3],
			Ssid:     data[2],
			Mac:      data[0],
			Status:   "Cracked",
		}
		return &response, nil

	} else {
		response := models.CrackResult{
			Status: "Exhausted",
		}
		return &response, nil
	}
}

func exitStatus(state *os.ProcessState) int {
	status, ok := state.Sys().(syscall.WaitStatus)
	if !ok {
		return -1
	}
	return status.ExitStatus()
}
