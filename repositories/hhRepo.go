package repositories

import "database/sql"

type HandshakeRepo struct {
	db *sql.DB
}

func NewHandshakeRepo(db *sql.DB) *HandshakeRepo {
	return &HandshakeRepo{
		db: db,
	}
}

func (r *HandshakeRepo) FindByID(ID int) {
	//dsafdf
}
