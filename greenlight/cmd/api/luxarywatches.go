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

	watches := data.Watches{
		ID:          id,
		CreatedAt:   time.Now(),
		Title:       "Rolex",
		Runtime:     102,
		WatchesType: []string{"Submariner", "Daytona", "Datejust"},
		Version:     1,
	}

	// Create an envelope {"watchese": watches} instance and pass it to writeJSON(), instead
	// of passing the plain movie struct.
	err = app.writeJSON(w, http.StatusOK, envelope{"watches": watches}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
