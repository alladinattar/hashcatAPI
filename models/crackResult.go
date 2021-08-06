package models

import "os"

type CrackResult struct {
	Ssid     string `json:"ssid,omitempty"`
	Password string `json:"password,omitempty"`
	Mac      string `json:"mac,omitempty"`
	Status   string `json:"status"`
}

type Cracker interface {
	CrackWPA(file *os.File) (*CrackResult, error)
}
