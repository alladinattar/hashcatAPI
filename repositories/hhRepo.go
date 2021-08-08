package repositories

import (
	"database/sql"
	"fmt"
	"github.com/hashcatAPI/models"
	_ "github.com/mattn/go-sqlite3" //for sqlite database
	"log"
)

type HandshakeRepository struct {
	db *sql.DB
}

func NewHandshakeRepository(db *sql.DB) *HandshakeRepository {
	return &HandshakeRepository{
		db: db,
	}
}

func (r *HandshakeRepository) GetHandshakes() ([]*models.Handshake, error) {
	rows, err := r.db.Query("SELECT * FROM handshakes")
	if err != nil {
		log.Println(err)
	}
	var handshakes []*models.Handshake

	for rows.Next() {
		handshake := &models.Handshake{}
		err = rows.Scan(&handshake.ID, &handshake.MAC, &handshake.SSID, &handshake.Encryption,
			&handshake.Latitude, &handshake.Longitude, &handshake.IMEI, &handshake.Time, &handshake.Password)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(handshake)
		handshakes = append(handshakes, handshake)
	}

	rows.Close()
	return handshakes, err
}
func (r *HandshakeRepository) GetByID(ID int) (*models.Handshake, error) {
	return &models.Handshake{}, nil
}

func (r *HandshakeRepository) GetAll() ([]*models.Handshake, error) {
	return nil, nil
}
func (r *HandshakeRepository) Save(handshake *models.Handshake) (int, error) {
	return 0, nil
}

func (r *HandshakeRepository) GetByMAC(MAC string) (*models.Handshake, error) {
	return &models.Handshake{}, nil
}
