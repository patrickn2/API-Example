package database

import (
	"log"

	"gorm.io/gorm"
)

func InitDB(openDB gorm.Dialector, config *gorm.Config) *gorm.DB {
	db, err := gorm.Open(openDB, config)
	if err != nil {
		log.Fatalln("error connecting to database ", err)
	}
	return db
}
