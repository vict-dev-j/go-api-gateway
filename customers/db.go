package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	//_ "github.com/golang-migrate/migrate/v4/database/postgres"
	//_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "customers")

	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
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
	err = runMigrations(dbInfo)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

//func runMigrations(dbInfo string) error {
//	m, err := migrate.New("customers/migrations", dbInfo)
//	if err != nil {
//		return err
//	}
//
//	err = m.Up()
//	if err != nil && err != migrate.ErrNoChange {
//		return err
//	}
//
//	log.Println("Database migration successful")
//	return nil
//}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
