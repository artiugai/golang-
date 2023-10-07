package data

import (
	"time"
)

// Watches represents information about a watches.
type Watches struct {
	ID          int64     // Unique integer ID for the watches
	CreatedAt   time.Time // Timestamp for when the watches is added to our database
	Title       string    // Watches title
	Year        int32     // Watches release year
	Runtime     int32     //  --
	WatchesType []string  // Type of watches
	Version     int32     // The version number starts at 1 and will be incremented each time the movie information is updated
}
