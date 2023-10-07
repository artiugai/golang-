package main

import (
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"net/http"
	"time"
)

func (app *application) createWatchesHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "create a new watches")
}

func (app *application) showWatchesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Create a new instance of the Watches struct, containing the ID we extracted from
	// the URL and some dummy data. Also notice that we deliberately haven't set a
	// value for the Year field.
	watches := data.Watches{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Rolex",
		Runtime:     102,
		WatchesType: []string{"Datejust", "Submariner", "Daytona"},
		Version:     1,
	}

	// Encode the struct to JSON and send it as the HTTP response.
	err = app.writeJSON(w, http.StatusOK, watches, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
