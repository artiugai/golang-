package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/watches", app.listWatchesHandler)          // Изменено с accessories на watches
	router.HandlerFunc(http.MethodPost, "/v1/watches", app.createWatchesHandler)       // Изменено с accessories на watches
	router.HandlerFunc(http.MethodGet, "/v1/watches/:id", app.showWatchesHandler)      // Изменено с accessories на watches
	router.HandlerFunc(http.MethodPatch, "/v1/watches/:id", app.updateWatchesHandler)  // Изменено с accessories на watches
	router.HandlerFunc(http.MethodDelete, "/v1/watches/:id", app.deleteWatchesHandler) // Изменено с accessories на watches

	return router
}
