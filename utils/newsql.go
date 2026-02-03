package utils

import (
	"database/sql"
	"github.com/spf13/viper"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", viper.GetString("database.path"))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
