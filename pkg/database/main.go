package database

import (
	"os"

	"github.com/jinzhu/gorm"
)

// NewGormConnect hoge
func NewGormConnect() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err.Error())
	}
	return db
}
