package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the WatchesModel and UserModel.
type Models struct {
	Watches WatchesModel
	Users   UserModel // Add a new Users field.
}

func NewModels(db *sql.DB) Models {
	return Models{
		Watches: WatchesModel{DB: db},
		Users:   UserModel{DB: db}, // Initialize a new UserModel instance.
	}
}
