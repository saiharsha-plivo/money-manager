package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/saiharsha/money-manager/internal/data"
	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
	"github.com/saiharsha/money-manager/pkg/validator"
)

func (app *application) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := app.models.Currencies.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	envelope := jsonhelper.Envelope{
		"currencies": currencies,
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) InsertCurrency(w http.ResponseWriter, r *http.Request) {
	var currency data.Currency

	err := jsonhelper.ReadJSON(w, r, &currency)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.logger.PrintDebug("properties of currency", map[string]string{
		"currency-name": currency.Name,
		"cuurency-rate": fmt.Sprintf("%v", currency.Rate),
	})

	v := validator.NewValidator()

	err = app.models.Currencies.Insert(&currency)
	if err != nil {
		app.logger.PrintDebug("error while inserting data", nil)
		if errors.Is(err, data.ErrDuplicateCurrency) {
			v.AddError("currency", "a currency with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	envelope := jsonhelper.Envelope{
		"currency": currency,
	}

	err = jsonhelper.WriteJSON(w, http.StatusCreated, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Currencies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, 200, jsonhelper.Envelope{"message": "currency successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetCurrency(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currency, err := app.models.Currencies.GetByID(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.PrintDebug("message after getting currency by id", map[string]string{
		"currency-name": currency.Name,
		"currency-rate": fmt.Sprintf("%v", currency.Rate),
	})

	err = jsonhelper.WriteJSON(w, http.StatusAccepted, jsonhelper.Envelope{"currency": currency}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	currency, err := app.models.Currencies.GetByID(id)

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
		Name *string  `json:"name"`
		Rate *float32 `json:"rate"`
	}

	err = jsonhelper.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		currency.Name = *input.Name
	}

	if input.Rate != nil {
		currency.Rate = *input.Rate
	}

	app.logger.PrintDebug("currency after update", map[string]string{
		"currency-name":    currency.Name,
		"currency-rate":    fmt.Sprintf("%v", currency.Rate),
		"currency-version": fmt.Sprintf("%v", currency.Version),
	})

	err = app.models.Currencies.Update(currency)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"currency": currency}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
