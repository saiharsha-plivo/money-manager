package main

import (
	"errors"
	"net/http"

	"github.com/saiharsha/money-manager/internal/data"
	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
	"github.com/saiharsha/money-manager/pkg/validator"
)

func (app *application) GetAllRecordTypes(w http.ResponseWriter, r *http.Request) {
	recordtypes, err := app.models.RecordTypes.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	envelope := jsonhelper.Envelope{
		"recordtypes": recordtypes,
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) InsertRecordType(w http.ResponseWriter, r *http.Request) {
	var recordtype data.RecordType

	err := jsonhelper.ReadJSON(w, r, &recordtype)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()

	err = app.models.RecordTypes.Insert(&recordtype)
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
		"recordtype": recordtype,
	}

	err = jsonhelper.WriteJSON(w, http.StatusCreated, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeleteRecordType(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.RecordTypes.Delete(id)
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

func (app *application) GetRecordType(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	recordtype, err := app.models.RecordTypes.GetByID(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusAccepted, jsonhelper.Envelope{"recordtype": recordtype}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
