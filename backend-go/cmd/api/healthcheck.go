package main

import (
	"net/http"

	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	env := jsonhelper.Envelope{
		"status":      "available",
		"environment": app.config.environment,
		"version":     "1.0.0",
	}
	err := jsonhelper.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
