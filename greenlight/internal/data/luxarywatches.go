package data

import "time"

type Watches struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	// Use the Runtime type instead of int32.
	// Note: The omitempty directive will still work on this: if the Runtime field has the
	// underlying value 0, then it will be considered empty and omitted.
	Runtime     Runtime  `json:"runtime,omitempty"`
	WatchesType []string `json:"WatchesType,omitempty"`
	Version     int32    `json:"version"`
}
