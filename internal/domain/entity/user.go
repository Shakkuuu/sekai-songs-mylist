package entity

import "time"

type User struct {
	ID             string
	Email          string
	Password       string
	IsVerified     bool
	VerifyToken    string
	TokenExpiresAt time.Time
	IsAdmin        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
