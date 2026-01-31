package utils

import (
	"database/sql"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./history.db")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
