package repositories

import (
	"database/sql"
	"errors"
	"github.com/hashcatAPI/models"
	_ "github.com/mattn/go-sqlite3" //for sqlite database
	"log"
	"strings"
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
		err = rows.Scan(&handshake.ID, &handshake.MAC, &handshake.SSID, &handshake.Encryption,
			/*&handshake.Latitude, &handshake.Longitude, &handshake.IMEI,*/ &handshake.Time, &handshake.Password)
		if err != nil {
			log.Println(err)
		}
		handshakes = append(handshakes, handshake)
	}
	rows.Close()
	return handshakes, err
}

func (r *HandshakeRepository) Save(handshakes []*models.Handshake) (int, error) {
	originalHandshakes, repeatedHandshakes := r.checkHandshakes(handshakes)
	if len(repeatedHandshakes) != 0 {
		return 0, errors.New("These handshakes already exists: " + strings.Join(repeatedHandshakes, ",") + ". Others added successfully")
	}
	for _, handshake := range originalHandshakes {
		stmt, err := r.db.Prepare("INSERT INTO handshakes(mac , ssid , password, time,enctyption) values(?,?,?,?,?)")
		if err != nil {
			log.Println("Failed prepare insert query", err)
			return 1, err
		}
		_, err = stmt.Exec(handshake.MAC, handshake.SSID, handshake.Password, handshake.Time, handshake.Encryption)
		if err != nil {
			log.Println("Failed exec insert query", err)
			return 1, err
		}

	}
	return len(originalHandshakes), nil
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

func (r *HandshakeRepository) checkHandshakes(handshakes []*models.Handshake) (originalHandshakes []*models.Handshake, repeatedHandshakes []string) {
	for _, handshake := range handshakes {
		if handshake.SSID == "" || handshake.Password == "" || handshake.MAC == "" {
			continue
		}
		handshakes, err := r.GetByMAC(handshake.MAC)
		if err != nil {
			return nil, repeatedHandshakes
		}
		if len(handshakes) != 0 {
			repeatedHandshakes = append(repeatedHandshakes, handshake.MAC)
			continue
		}
		originalHandshakes = append(originalHandshakes, handshake)
	}
	return originalHandshakes, repeatedHandshakes
}
