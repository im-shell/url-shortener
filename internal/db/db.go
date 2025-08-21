package db

// create postgres connection

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Create() (*sql.DB, error) {
	host:= os.Getenv("DB_HOST")
	port:= os.Getenv("DB_PORT")
	user:= os.Getenv("DB_USER")
	password:= os.Getenv("DB_PASSWORD")
	dbname:= os.Getenv("DB_NAME")

	connStr:= fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err:= sql.Open("postgres", connStr)
	if err != nil {
		return nil , fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	return db, nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}