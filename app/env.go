package app

import (
	"database/sql"
)

type Env struct {
	db *sql.DB
}

func NewEnv(db *sql.DB) *Env {
	return &Env{db}
}
