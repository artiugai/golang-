package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when // looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this, // like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	// Movies field is an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Watches interface {
		Insert(watches *Watches) error
		Get(id int64) (*Watches, error)
		Update(watches *Watches) error
		Delete(id int64) error
	}
}

// For ease of use, we also add a New() method which returns a Models struct containing // the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Watches: WatchesModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Watches: MockWatchesModel{},
	}
}
