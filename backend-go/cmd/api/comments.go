package main

import (
	"errors"
	"net/http"

	"github.com/saiharsha/money-manager/internal/data"
	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
	"github.com/saiharsha/money-manager/pkg/validator"
)

func (app *application) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var v = validator.NewValidator()
	var filters data.Filters
	filters.Page = jsonhelper.ReadIntParam(r, "page", 1, v)
	filters.PageSize = jsonhelper.ReadIntParam(r, "page_size", 20, v)
	filters.Sort = jsonhelper.ReadStringParam(r, "sort", "id")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	comments, err := app.models.Comments.GetAll(id, filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"comments": comments}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment data.Comment
	err := jsonhelper.ReadJSON(w, r, &comment)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateComment(v, &comment)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Comments.Insert(&comment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateComment):
			v.AddError("comment", "a comment with this description already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusCreated, jsonhelper.Envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Description *string `json:"description"`
		RecordID    *int64  `json:"record_id"`
	}
	err = jsonhelper.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	data.ValidateComment(v, &data.Comment{
		Description: *input.Description,
		RecordID:    *input.RecordID,
	})

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	comment, err := app.models.Comments.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Description != nil {
		comment.Description = *input.Description
	}

	if input.RecordID != nil {
		comment.RecordID = *input.RecordID
	}

	err = app.models.Comments.Update(comment)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := jsonhelper.ReadIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Comments.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"message": "comment deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
