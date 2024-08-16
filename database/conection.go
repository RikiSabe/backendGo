package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const dsn = "host=181.188.156.195 user=xxxusxrdialogo password=tarija2024 dbname=testRS port=18020"

func Conection() {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal((err))
	} else {
		log.Println("Database Conected")
	}
}
