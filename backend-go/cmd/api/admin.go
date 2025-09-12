package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/saiharsha/money-manager/internal/data"
	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
	"github.com/saiharsha/money-manager/pkg/validator"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleSuperUser Role = "superuser"
)

func (app *application) UpdateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Role  Role   `json:"role"`
	}

	err := jsonhelper.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.NewValidator()
	v.Check(input.Role == "user" || input.Role == "superuser" || input.Role == "admin", "user", fmt.Sprintf("role type %v not allowed", input.Role))

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetUserByMail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "user with email not found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.PrintDebug("got the userdetails from email", map[string]string{
		"user": user.Name,
		"role": user.Role,
	})

	user.Role = string(input.Role)
	updateduser, err := app.models.Users.UpdateUser(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.PrintDebug("patched the user for update role succesfully", map[string]string{
		"user": updateduser.Name,
		"role": updateduser.Role,
	})

	err = jsonhelper.WriteJSON(w, http.StatusAccepted, jsonhelper.Envelope{"message": "user patched succesfully", "updateduser": updateduser}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
