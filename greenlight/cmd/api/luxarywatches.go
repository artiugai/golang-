package main

import (
	"encoding/json"
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"net/http"
	"time"
)

func (app *application) createWatchesHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information expected in the HTTP request body.
	var input struct {
		Title       string   `json:"title"`
		Year        int32    `json:"year"`
		Runtime     int32    `json:"runtime"`
		WatchesType []string `json:"watchesType"`
	}

	// Initialize a new json.Decoder instance which reads from the request body,
	// and then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode(), we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Dump the contents of the input struct in an HTTP response.
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
