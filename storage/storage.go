package storage

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AnyoneClown/CocaCallsAPI/types"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CockroachDB struct {
	db *gorm.DB
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
		log.Fatalf("COCKROACH_DB_URL environment variable not set")
	}

	
	// Auto migrate models
	db.AutoMigrate(&types.User{})

	return &CockroachDB{db: db}
}

func (c *CockroachDB) CreateUser(email, password string) (types.User, error) {
	if err := types.ValidateUser(email, password); err != nil {
		return types.User{}, err
	}

	var existingUser types.User
    if err := c.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
        return types.User{}, fmt.Errorf("email already in use")
    } else if err != gorm.ErrRecordNotFound {
        return types.User{}, err
    }

	user := types.User{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := c.db.Create(&user)
	if result.Error != nil {
		return types.User{}, result.Error
	}

	return user, nil
}

func (c *CockroachDB) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	result := c.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, result.Error
	}

	return user, nil
}
