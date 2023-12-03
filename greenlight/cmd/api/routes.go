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
	router.HandlerFunc(http.MethodGet, "/v1/watches", app.listWatchesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/watches", app.createWatchesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/watches/:id", app.showWatchesHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/watches/:id", app.updateWatchesHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/watches/:id", app.deleteWatchesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
