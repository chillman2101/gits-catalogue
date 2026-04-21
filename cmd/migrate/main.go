package main

import (
	"log"

	"github.com/chillman2101/gits-catalogue/internal/config"
	"github.com/chillman2101/gits-catalogue/internal/model"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	cfg := config.Load()
	// create database if not exist
	if err := cfg.CreateDB(); err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	db, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("running migrations...")

	if err := db.AutoMigrate(
		&model.User{},
		&model.Author{},
		&model.Publisher{},
		&model.Book{},
	); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	log.Println("migrations completed successfully")
}
