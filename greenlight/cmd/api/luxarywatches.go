package main

import (
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/validator"
	"net/http"
	"time"
)

func (app *application) createWatchesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string     `json:"title"`
		Year  int32      `json:"year,omitempty"`
		Price data.Price `json:"price,omitempty"`
		Type  []string   `json:"type,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Создаем переменную watches, содержащую указатель на структуру Watches.
	watches := &data.Watches{
		Title: input.Title,
		Year:  input.Year,
		Price: input.Price,
		Type:  input.Type,
	}

	v := validator.New()
	if data.ValidateWatches(v, watches); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Watches.Insert(watches)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/watches/%d", watches.ID))

	// Отправляем JSON-ответ со статус-кодом 201 Created, данными о часах в теле ответа
	// и заголовком Location.
	err = app.writeJSON(w, http.StatusCreated, envelope{"watches": watches}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(input.Price != 0, "price", "must be provided")
	v.Check(input.Price > 0, "price", "must be a positive float")
	v.Check(input.Type != nil, "watchesType", "must be provided")
	v.Check(len(input.Type) >= 1, "watchesType", "must contain at least 1 watches type")
	v.Check(len(input.Type) <= 5, "watchesType", "must not contain more than 5 types")
	// Note that we're using the Unique helper in the line below to check that all
	// values in the input.Genres slice are unique. v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
	// Use the Valid() method to see if any of the checks failed. If they did, then use // the failedValidationResponse() helper to send a response to the client, passing // in the v.Errors map.
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Rolex",
		Price:     50.000,
		Type:      []string{"Submariner"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"watches": watches}, nil)
	if err != nil {
		// Use the new serverErrorResponse() helper.
		app.serverErrorResponse(w, r, err)
	}
}
