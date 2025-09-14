package main

import (
	"fmt"
	"net/http"
)

func (app *application) panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.errorResponse(w, r, http.StatusInternalServerError, fmt.Sprintf("the server encountered a problem and could not process your request: %v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) BackgroundTask(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.wg.Add(1)
		go func() {
			defer app.wg.Done()
			next.ServeHTTP(w, r)
		}()
	})
}

func (app *application) BackgroundEmailTask(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		fn()
	}()
}
