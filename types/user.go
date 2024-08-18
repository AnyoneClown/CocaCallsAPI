package types

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"uuid"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}
