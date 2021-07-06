package repositories

import (
	"database/sql"
	"github.com/hashcatAPI/models"
)

type HandshakeRepo struct {
	db *sql.DB
}

func NewHandshakeRepo(db *sql.DB) *HandshakeRepo {
	return &HandshakeRepo{
		db: db,
	}
}

func (r *HandshakeRepo) GetByID(ID int) (models.Handshake, error) {
	//dsafdf
	return models.Handshake{}, nil
}

func (r *HandshakeRepo) Save() (int, error) {
	return 0, nil
}

func (r *HandshakeRepo) GetByMAC(MAC string) (models.Handshake, error) {

}
