package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (app *application) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		message := fmt.Sprintf("recieved %v request from %v at %v", r.Method, r.RemoteAddr, r.URL.String())
		app.logger.PrintInfo(message, map[string]string{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		})
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		message = fmt.Sprintf("processed %v request from %v at %v", r.Method, r.RemoteAddr, r.URL.String())
		app.logger.PrintInfo(message, map[string]string{
			"method":   r.Method,
			"url":      r.URL.String(),
			"ip":       r.RemoteAddr,
			"duration": duration.String(),
		})
	})
}

func (app *application) VerifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "

		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		app.logger.PrintDebug("auth header", map[string]string{"raw": authHeader})

		if authHeader == "" || !strings.HasPrefix(authHeader, prefix) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
		if token == "" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.VerifyToken(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		app.logger.PrintDebug(fmt.Sprintf("user from jwt decode : %v", user), nil)
		r = app.contextSetUser(r, user)

		app.logger.PrintDebug(fmt.Sprintf("context set in request: %v", app.contextGetUser(r)), nil)
		app.logger.PrintDebug(fmt.Sprintf("request after adding context: %v", *r), nil)

		next.ServeHTTP(w, r)
	})
}

func (app *application) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := app.contextGetUser(r)
			if user == nil {
				app.invalidAuthenticationTokenResponse(w, r)
				return
			}

			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			app.notPermittedResponse(w, r)
		})
	}
}
