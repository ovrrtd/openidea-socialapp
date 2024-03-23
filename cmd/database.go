package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func newDBDefaultSql() (*sql.DB, error) {
	dsn := ""
	fmt.Println("env :", os.Getenv("ENV"))
	if os.Getenv("ENV") != "production" {
		// Create a PostgreSQL database connection string
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
			DB_HOST,
			DB_USERNAME,
			DB_PASSWORD,
			DB_NAME,
			DB_PORT,
		)
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=verify-full sslrootcert=ap-southeast-1-bundle.pem TimeZone=UTC",
			DB_HOST,
			DB_USERNAME,
			DB_PASSWORD,
			DB_NAME,
			DB_PORT,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return db, nil
}
