package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/saiharsha/money-manager/internal/data"
	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
	"github.com/saiharsha/money-manager/pkg/validator"
)

func (app *application) listRecordsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var v = validator.NewValidator()
	var filters data.Filters
	filters.Page = jsonhelper.ReadIntParam(r, "page", 1, v)
	filters.PageSize = jsonhelper.ReadIntParam(r, "page_size", 20, v)
	filters.Sort = jsonhelper.ReadStringParam(r, "sort", "id")

	// no default value for start date
	filters.StartDate = jsonhelper.ReadTimeParam(r, "start_date", time.Time{}, false, v)
	filters.EndDate = jsonhelper.ReadTimeParam(r, "end_date", time.Now(), true, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	records, err := app.models.Records.GetRecordsForUser(userID, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"metadata": filters.CalculateMetadata(len(records)), "records": records}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	var record data.Record

	err := jsonhelper.ReadJSON(w, r, &record)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateRecord(v, &record)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Records.Insert(&record)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateRecord):
			v.AddError("record", "record already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusCreated, jsonhelper.Envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.models.Records.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.models.Records.GetByID(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Amount      *int64  `json:"amount"`
		Description *string `json:"description"`
		TypeID      *int64  `json:"type_id"`
		CurrencyID  *int64  `json:"currency_id"`
	}

	err = jsonhelper.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateRecord(v, &data.Record{
		Amount:      *input.Amount,
		Description: *input.Description,
		TypeID:      *input.TypeID,
		CurrencyID:  *input.CurrencyID,
	})

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if input.Amount != nil {
		record.Amount = *input.Amount
	}
	if input.Description != nil {
		record.Description = *input.Description
	}
	if input.TypeID != nil {
		record.TypeID = *input.TypeID
	}
	if input.CurrencyID != nil {
		record.CurrencyID = *input.CurrencyID
	}

	err = app.models.Records.Update(record)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Records.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"message": "record successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
