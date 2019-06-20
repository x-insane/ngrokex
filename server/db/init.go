package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func GetConnection() (*gorm.DB, error) {
	client, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		return nil, err
	}
	return client, nil
}