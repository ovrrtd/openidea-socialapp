package main

import (
	"database/sql"
	"fmt"
	"time"
)

func newDBDefaultSql() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=UTC %s",
		DB_HOST,
		DB_USERNAME,
		DB_PASSWORD,
		DB_NAME,
		DB_PORT,
		DB_PARAMS,
	)

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
