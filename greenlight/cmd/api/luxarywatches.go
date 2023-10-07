package main

import (
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"net/http"
	"time"
)

func (app *application) createWatchesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Year        int32    `json:"year"`
		Runtime     int32    `json:"runtime"`
		WatchesType []string `json:"watchesType"`
	}
	// Use the new readJSON() helper to decode the request body into the input struct.
	// If this returns an error we send the client the error message along with a 400
	// Bad Request status code, just like before.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showWatchesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"watches": watches}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
