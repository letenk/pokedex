package util

import (
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// counts for count time response database
var counts int64

// Function open connection to Database
func openDB(dsn string) (*gorm.DB, error) {
	// Open connection to db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Function for connection to db
func SetupDB(dsn string) *gorm.DB {
	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Fatalf("Database Connection error: %s\n", err)
			return nil
		}

		log.Println("Backing off for two seconds ...")
		time.Sleep(2 * time.Second)
		continue
	}
}
