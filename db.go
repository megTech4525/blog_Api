package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sqlx.Connect("pgx", dsn)

	if err != nil {
		log.Println(err)

	}

	if err := db.Ping(); err != nil {
		log.Println(err)

	}

	db.SetMaxIdleConns(25)
	db.SetMaxIdleConns(25)
	log.Println("db is connected", os.Getenv("DB_NAME"))
	return db

}
