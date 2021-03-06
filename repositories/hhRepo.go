package repositories

import (
	"database/sql"
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

func (r *HandshakeRepository) GetByID(ID int) (*models.Handshake, error) {
	return &models.Handshake{}, nil
}

func (r *HandshakeRepository) GetAll() (handshakes []*models.Handshake, err error) {
	rows, err := r.db.Query("SELECT * FROM handshakes")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		handshake := &models.Handshake{}
		err = rows.Scan(&handshake.ID, &handshake.MAC, &handshake.SSID, &handshake.Password, &handshake.Time,
			&handshake.Encryption, &handshake.Longitude, &handshake.Latitude, &handshake.IMEI)
		if err != nil {
			log.Println(err)
		}
		handshakes = append(handshakes, handshake)
	}
	rows.Close()
	return handshakes, err
}

func (r *HandshakeRepository) Save(handshake *models.Handshake) (int, error) {
	affectedRows := 0
	stmt, err := r.db.Prepare("INSERT INTO handshakes(mac , ssid , password, time,enctyption, longitude, latitude, imei) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("Failed prepare insert query", err)
		return 0, nil
	}

	result, err := stmt.Exec(handshake.MAC, handshake.SSID, handshake.Password, handshake.Time, handshake.Encryption, handshake.Longitude, handshake.Latitude, handshake.IMEI)
	if err != nil {
		log.Println("Failed exec insert query", err)
		return 0, nil
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed get rows affected", err)
	}
	affectedRows += int(rowsAffected)
	return affectedRows, nil
}

func (r *HandshakeRepository) GetByMAC(MAC string) (handshakes []*models.Handshake, err error) {
	rows, err := r.db.Query("SELECT mac FROM handshakes WHERE mac='" + MAC + "'")
	if err != nil {
		return nil, err
	}
	handshake := models.Handshake{}
	for rows.Next() {
		err = rows.Scan(&handshake.MAC)
		if err != nil {
			return nil, err
		}
		handshakes = append(handshakes, &handshake)
	}
	rows.Close()
	return handshakes, nil
}
