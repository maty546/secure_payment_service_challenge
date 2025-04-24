package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	//todo make this a config
	dsn := "host=localhost user=postgres password=ciclon16 dbname=postgres port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to database:", err)
	}
	fmt.Println("✅ Connected to PostgreSQL!")
	return database

}
