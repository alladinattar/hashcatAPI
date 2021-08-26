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
	rows, err := r.db.Query("SELECT mac, ssid, password, imei, file FROM handshakes")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		handshake := &models.Handshake{}
		err = rows.Scan(&handshake.MAC, &handshake.SSID, &handshake.Password, &handshake.IMEI, &handshake.File)
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

func (r *HandshakeRepository) AddTaskToDB(task *models.Handshake) error {
	stmt, err := r.db.Prepare("INSERT INTO tasks(filename, imei, status) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(task.File, task.IMEI, task.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *HandshakeRepository) UpdateTaskState(task *models.Handshake) error {
	stmt, err := r.db.Prepare("update tasks set status=? where imei=? and filename=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(task.Status, task.IMEI, task.File)
	if err != nil {
		return err
	}
	return nil
}

func (r *HandshakeRepository) GetFilesByIMEI(imei string)(files []string, err error){
	rows, err := r.db.Query("SELECT file FROM handshakes WHERE imei='" + imei + "'")
	if err != nil {
		return  nil, err
	}
	var file string
	for rows.Next() {
		err = rows.Scan(&file)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	rows.Close()
	return files, nil
}


func(r *HandshakeRepository) GetProgressByIMEI(imei string)(files []*models.Handshake, err error){
	rows, err := r.db.Query("SELECT filename, status FROM tasks where imei = '" + imei + "'")
	if err != nil {
		return  nil, err
	}
	for rows.Next() {
		var handshake models.Handshake
		err = rows.Scan(&handshake.File, &handshake.Status)
		if err != nil {
			return nil, err
		}
		files = append(files, &handshake)
	}
	rows.Close()
	return files, nil
}