package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *application) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(app.RequestLogger)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	// health check
	r.Get("/", app.healthcheck)

	// account
	r.Route("/users", func(r chi.Router) {
		r.Post("/signup", app.UserSignUp)
		r.Post("/login", app.UserLogin)
		r.Post("/signout", app.UserSignOut)
		r.Post("/refresh", app.UserRefresh)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.VerifyUser)
		r.Route("/admin", func(r chi.Router) {
			r.Use(app.RequireRole("admin"))
			r.Patch("/role", app.ChangeUserRole)
			r.Get("/currencies", app.GetAllCurrencies)
			r.Post("/currencies", app.InsertCurrency)
			r.Get("/currencies/{id}", app.GetCurrency)
			r.Patch("/currencies/{id}", app.UpdateCurrency)
			r.Delete("/currencies/{id}", app.DeleteCurrency)
		})
	})

	return r
}
