package main

import (
	"time"

	"github.com/alinaqigit/rss-generator-project/internal/db"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username      string    `json:"name"`
}

func database_user_to_User(databaseUser db.User) User {
	return User{
		ID: databaseUser.ID,
		CreatedAt: databaseUser.CreatedAt,
    UpdatedAt: databaseUser.UpdatedAt,
    Username: databaseUser.Username,
	}
}