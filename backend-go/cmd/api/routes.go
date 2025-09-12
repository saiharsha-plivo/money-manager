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

	// records
	r.Route("/records", func(r chi.Router) {
		r.Use(app.VerifyUser)
		r.Get("/", app.listRecordsHandler)
		r.Post("/", app.createRecordHandler)
		r.Get("/{id}", app.getRecordHandler)
		r.Patch("/{id}", app.updateRecordHandler)
		r.Delete("/{id}", app.deleteRecordHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(app.VerifyUser)
		r.Route("/admin", func(r chi.Router) {
			r.Use(app.RequireRole("admin"))
			r.Patch("/role", app.UpdateUserRoleHandler)
			r.Get("/currencies", app.GetCurrenciesHandler)
			r.Post("/currencies", app.InsertCurrencyHandler)
			r.Get("/currencies/{id}", app.GetCurrencyHandler)
			r.Patch("/currencies/{id}", app.UpdateCurrencyHandler)
			r.Delete("/currencies/{id}", app.DeleteCurrencyHandler)

			r.Get("/recordtypes", app.GetRecordTypesHandler)
			r.Post("/recordtypes", app.InsertRecordTypeHandler)
			r.Get("/recordtypes/{id}", app.GetRecordTypeHandler)
			r.Delete("/recordtypes/{id}", app.DeleteRecordTypeHandler)
		})
	})

	return r
}
