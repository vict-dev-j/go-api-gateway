package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func initDB() {
	var err error

	dbConfig := DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	maxAttempts := 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database")
			break
		}
		log.Printf("Error pinging the database (attempt %d/%d): %s", attempts, maxAttempts, err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to the database after %d attempts", maxAttempts)
	}
}
