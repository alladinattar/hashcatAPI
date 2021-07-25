package models

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
}

type HandshakeRepository interface {
	Save(string, string, string) (int, error)
	GetByID(ID int) (*Handshake, error)
	GetByMAC(MAC string) (*Handshake, error)
	GetHandshakes() ([]*Handshake, error)
}
