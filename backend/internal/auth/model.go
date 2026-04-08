package auth

import "time"

type User struct {
	ID              string    `db:"id"`
	EstablishmentID string    `db:"establishment_id"`
	Name            string    `db:"name"`
	Email           string    `db:"email"`
	PasswordHash    string    `db:"password_hash"`
	Role            string    `db:"role"`
	Active          bool      `db:"active"`
	CreatedAt       time.Time `db:"created_at"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RegisterInput struct {
	OwnerName         string
	EstablishmentName string
	Email             string
	Password          string
	Slug              string
	ContactPhone      *string
}
