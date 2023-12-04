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

type Models struct {
	Watches     WatchesModel
	Tokens      TokenModel // Add a new Tokens field.
	Permissions PermissionModel
	Users       UserModel
}

func NewWatchesModel(db *sql.DB) Models {
	return Models{
		Watches:     WatchesModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db}, // Initialize a new TokenModel instance.
		Users:       UserModel{DB: db},
	}
}
