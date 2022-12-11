package util

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// counts for count time response database
var counts int64

// Function open connection to Database
func openDB(dsn string) (*sql.DB, error) {
	// Open connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Ping db
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Function for connection to db
func SetupDB(dsn string) *sql.DB {
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
