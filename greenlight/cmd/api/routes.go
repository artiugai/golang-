package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listWatchesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createWatchesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showWatchesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateWatchesHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteWatchesHandler)

	// Wrap the router with the panic recovery middleware.
	return app.recoverPanic(app.rateLimit(router))
}
