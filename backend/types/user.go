package types

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Email     string    `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserWithoutPassword struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Validate() error {
	// Validate pattern for email
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Validate password (minimum length of 8 characters)
	if len(u.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	return nil
}

func ValidateUser(email, password string) error {
	user := User{
		Email:    email,
		Password: password,
	}

	if err := user.Validate(); err != nil {
		return err
	}

	return nil
}