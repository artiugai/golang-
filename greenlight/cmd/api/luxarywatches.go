package main

import (
	"errors"
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/validator"
	"net/http"
)

func (app *application) createWatchesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string   `json:"title"`
		Year     int32    `json:"year"`
		Price    float64  `json:"price"`
		Brand    []string `json:"Brand"`
		Material []string `json:"material"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	watches := &data.Watches{
		Title:    input.Title,
		Year:     input.Year,
		Price:    input.Price,
		Brand:    input.Brand,
		Material: input.Material,
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
	err = app.writeJSON(w, http.StatusCreated, envelope{"watches": watches}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}
func (app *application) showWatchesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	watches, err := app.models.Watches.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"watches": watches}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateWatchesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	watches, err := app.models.Watches.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title    *string   `json:"title"`
		Year     *int32    `json:"year"`
		Price    *float64  `json:"price"`
		Brand    *[]string `json:"Brand"`
		Material *[]string `json:"Material"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		watches.Title = *input.Title
	}
	if input.Year != nil {
		watches.Year = *input.Year
	}
	if input.Price != nil {
		watches.Price = *input.Price
	}
	if input.Brand != nil {
		watches.Brand = *input.Brand
	}
	if input.Material != nil {
		watches.Material = *input.Material
	}

	v := validator.New()
	if data.ValidateWatches(v, watches); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Watches.Update(watches)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"watches": watches}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteWatchesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Watches.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "watches successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
