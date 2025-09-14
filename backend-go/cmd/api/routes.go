package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *application) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(app.RequestLogger)
	r.Use(app.panicRecover)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/", app.healthcheck)

	// account
	r.Route("/users", func(r chi.Router) {
		r.Post("/signup", app.UserSignUp)
		r.Post("/login", app.UserLogin)
		r.Post("/signout", app.UserSignOut)
		r.Post("/refresh", app.UserRefresh)
		r.Get("/verify/{token}", app.UserVerify)
	})

	// records
	r.Route("/records", func(r chi.Router) {
		r.Use(app.VerifyUser)
		r.Get("/list/{id}", app.listRecordsHandler)
		r.Post("/", app.createRecordHandler)
		r.Get("/{id}", app.getRecordHandler)
		r.Patch("/{id}", app.updateRecordHandler)
		r.Delete("/{id}", app.deleteRecordHandler)
	})

	// comments
	r.Route("/comments", func(r chi.Router) {
		r.Use(app.VerifyUser)
		r.Post("/", app.CreateCommentHandler)
		r.Get("/{id}", app.GetCommentsHandler)      // id is the record id
		r.Patch("/{id}", app.UpdateCommentHandler)  // id is the record id
		r.Delete("/{id}", app.DeleteCommentHandler) // id is comment id
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
