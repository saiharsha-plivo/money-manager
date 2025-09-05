package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/saiharsha/money-manager/internal/data"
	jsonhelper "github.com/saiharsha/money-manager/pkg/json"
	"github.com/saiharsha/money-manager/pkg/validator"
)

func (app *application) UserSignUp(w http.ResponseWriter, r *http.Request) {
	var userdetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	err := jsonhelper.ReadJSON(w, r, &userdetails)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	user := &data.User{
		Name:      userdetails.Username,
		Email:     userdetails.Email,
		Activated: false,
	}

	err = user.Password.SetPasswordHash(userdetails.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	uservalidator := validator.NewValidator()
	data.ValidateUser(uservalidator, user)

	if !uservalidator.Valid() {
		app.failedValidationResponse(w, r, uservalidator.Errors)
		app.logger.PrintDebug(fmt.Sprintf("user %v has failed data validation", user.Name), nil)
		return
	}
	app.logger.PrintDebug(fmt.Sprintf("user %v has passed data validation", user.Name), nil)

	err = app.models.Users.CreateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			uservalidator.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, uservalidator.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		app.logger.PrintDebug(fmt.Sprintf("user %v has not been created in database error : %v", user.Name, err.Error()), nil)
		return
	}
	app.logger.PrintDebug(fmt.Sprintf("user %v has been created in database", user.Name), nil)

	err = jsonhelper.WriteJSON(w, http.StatusAccepted, jsonhelper.Envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UserLogin(w http.ResponseWriter, r *http.Request) {
	var userdetails struct {
		Email    string
		Password string
	}

	err := jsonhelper.ReadJSON(w, r, &userdetails)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	v := validator.NewValidator()
	data.ValidateEmail(v, userdetails.Email)
	data.CheckPassword(v, userdetails.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	app.logger.PrintDebug(fmt.Sprintf("email %v has passed intial validation", userdetails.Email), nil)

	user, err := app.models.Users.GetUserByMail(userdetails.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.usernotFound(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	matches, err := user.Password.Matches(userdetails.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !matches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	app.logger.PrintDebug(fmt.Sprintf("user %v has passed password validation", user.Name), nil)
	accesstoken, err := app.CreateToken(user, 1*time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.PrintDebug(fmt.Sprintf("user %v got the access token", user.Name), nil)

	refreshtoken, err := app.CreateToken(user, 10*24*time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.logger.PrintDebug(fmt.Sprintf("user %v got the refresh token", user.Name), nil)
	cookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    refreshtoken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int((10 * 24 * time.Hour).Seconds()),
	}
	http.SetCookie(w, &cookie)

	envelope := jsonhelper.Envelope{
		"accesstoken": accesstoken,
	}

	err = jsonhelper.WriteJSON(w, http.StatusAccepted, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) UserSignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshtoken",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	envelope := jsonhelper.Envelope{
		"message": "user signed-out",
	}

	err := jsonhelper.WriteJSON(w, http.StatusAccepted, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) UserRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshtoken")

	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	tokenString := cookie.Value

	user, err := app.VerifyToken(tokenString)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidAuthenticationToken):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	accesstoken, err := app.CreateToken(user, 1*time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = jsonhelper.WriteJSON(w, http.StatusOK, jsonhelper.Envelope{"user": user, "accesstoken": accesstoken}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
