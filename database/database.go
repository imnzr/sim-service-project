package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func InitDb(databaseURL string) *sql.DB {
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatalf("failed to open MYSQL database connection: %v", err)
	}

	// Koneksi ke database untuk memastikan semuanya OK
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to MYSQL database: %v", err)
	}

	log.Println("Connecting to database ...")
	// Delay
	time.Sleep(time.Second)
	fmt.Println("Connected to database successfully")

	// Pool Koneksi
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(25)

	return db
}
