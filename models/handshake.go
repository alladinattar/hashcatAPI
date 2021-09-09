package models

import "os"

type Handshake struct {
	ID         int    `json:"-"`
	MAC        string `json:"mac,omitempty"`
	SSID       string `json:"ssid,omitempty"`
	Encryption string `json:"encryption,omitempty"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
	IMEI       string `json:"imei,omitempty"`
	Time       string `json:"time,omitempty"`
	Password   string `json:"password,omitempty"`
	File       string `json:"file,omitempty"`
	Status     string `json:"status,omitempty"`
}

type HandshakeRepository interface {
	Save(*Handshake) (int, error)
	GetByMAC(MAC string) ([]*Handshake, error)
	GetAll() ([]*Handshake, error)
	AddTaskToDB(*Handshake) error
	UpdateTaskState(*Handshake) error
	GetFilesByIMEI(imei string)(files []string, err error)
	GetProgressByIMEI(imei string)(files []*Handshake, err error)
	UpdatePasswordByMAC(mac string, password string)error
	SaveOriginHandshake(handshake *Handshake)(error)
	GetAllCrackedHandshakesByIMEI(imei string)(handshakes []*Handshake, err error)
}

type Cracker interface {
	CrackWPA(file *os.File) ([]*Handshake, error)
}

type TasksQueue interface {
	AddTask([]byte) error
}
