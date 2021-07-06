package models

type Handshake struct {
	SourceMAC string
	hash      string
	ESSID     string
	BSSID     string
}

type HandshakeRepo interface {
	Save() (int, error)
	GetByID() (Handshake, error)
	GetByMAC() (Handshake, error)
}
