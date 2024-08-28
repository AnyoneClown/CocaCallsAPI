package storage

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CockroachDB struct {
	DB *gorm.DB
}

func NewCockroachDB() *CockroachDB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	dsn := os.Getenv("COCKROACH_DB_URL")
	if dsn == "" {
		log.Fatalf("COCKROACH_DB_URL environment variable not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &CockroachDB{DB: db}
}
