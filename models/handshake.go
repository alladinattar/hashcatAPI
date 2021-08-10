package models

import "os"

type Handshake struct {
	ID         int     `json:"-"`
	MAC        string  `json:"mac"`
	SSID       string  `json:"ssid"`
	Encryption string  `json:"encryption"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	IMEI       string  `json:"imei"`
	Time       string  `json:"time"`
	Password   string  `json:"password"`
	Status     string  `json:"status"`
}

type HandshakeRepository interface {
	Save([]*Handshake) (int, error)
	GetByID(ID int) (*Handshake, error)
	GetByMAC(MAC string) ([]*Handshake, error)
	GetAll() ([]*Handshake, error)
}

type Cracker interface {
	CrackWPA(file *os.File) ([]*Handshake, error)
}
