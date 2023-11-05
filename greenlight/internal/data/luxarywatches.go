package data

import (
	"database/sql"
	"greenlight.alexedwards.net/internal/validator"
	"time"
)

type Watches struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	// Use the Runtime type instead of int32.
	// Note: The omitempty directive will still work on this: if the Runtime field has the
	// underlying value 0, then it will be considered empty and omitted.
	Price       Price    `json:"Price,omitempty"`
	WatchesType []string `json:"WatchesType,omitempty"`
	Version     int32    `json:"version"`
}

func ValidateMovie(v *validator.Validator, watches *Watches) {
	v.Check(watches.Title != "", "title", "must be provided")
	v.Check(len(watches.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(watches.Year != 0, "year", "must be provided")
	v.Check(watches.Year >= 1888, "year", "must be greater than 1888")
	v.Check(watches.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(watches.Price != 0, "price", "must be provided")
	v.Check(watches.Price > 0, "price", "must be a positive integer")
	v.Check(watches.WatchesType != nil, "watchestype", "must be provided")
	v.Check(len(watches.WatchesType) >= 1, "watchestype", "must contain at least 1 genre")
	v.Check(len(watches.WatchesType) <= 5, "watchestype", "must not contain more than 5 watchestype")
	v.Check(validator.Unique(watches.WatchesType), "watchestype", "must not contain duplicate values")

}

type WatchesModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m WatchesModel) Insert(movie *Watches) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m WatchesModel) Get(id int64) (*Watches, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m WatchesModel) Update(movie *Watches) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m WatchesModel) Delete(id int64) error {
	return nil
}

type MockWatchesModel struct{}

func (m MockWatchesModel) Insert(watches *Watches) error {

	return nil //
}

func (m MockWatchesModel) Get(id int64) (*Watches, error) {

	return nil, nil
}

func (m MockWatchesModel) Update(watches *Watches) error {

	return nil
}

func (m MockWatchesModel) Delete(id int64) error {

	return nil
}
