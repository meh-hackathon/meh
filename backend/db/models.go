package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Archived  bool      `json:"archived"`
}

type UserCredentials struct {
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at" db:"deleted_at"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
}

type QrCode struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OwnerID   uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
	Slug      string    `json:"slug" db:"slug"`
}

type Content struct {
	ID        uuid.UUID `json:"id" db:"id"`
	OwnerID   uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
	Type      string    `json:"type" db:"type"`
	Data      string    `json:"data" db:"data"`
}
