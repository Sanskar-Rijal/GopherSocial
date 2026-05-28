package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"social/internal/mailer"
	"social/internal/store"

	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UserWithToken struct {
	User *store.User
	Token string `json:"token"`
}

type registerUserHandlerResponse = SuccessResponse[store.User]

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	registerUserHandlerResponse
//	@Failure		400		{object}	ErrorResponseWrapper
//	@Failure		500		{object}	ErrorResponseWrapper
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	//has the password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()
	plainToken := uuid.New().String()

	//we will hash the token for email, but keep plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	//Store tthe user
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	//Send mail portion 

	userWithToken := &UserWithToken{
		User: user,
		Token: plainToken,
	}

	data := struct {
		Username string 
		ActivationUrl string 
	}{
		Username: user.Username,
		ActivationUrl: fmt.Sprintf("www.gophersocial.com/confirm/%s",plainToken),
	}


		status, err := app.mailer.Send(
			mailer.UserWelcomeTemplate,
			user.Username,
			user.Email,
			data,
			app.config.mail.isDevelopment,
		)

		if err != nil {

			  app.logger.Errorw("failed to send welcome email",
                "error", err,
                "email", user.Email,
            )

			//RollBack userCreation 
			if err := app.store.Users.Delete(ctx, user.ID); err != nil {				
				app.logger.Errorw("error deleting user", "error", err)
			}
			app.internalServerError(w,r,err)
			return 
		}
		
	// log success with status code
    app.logger.Infow("Email sent", "status code", status)
	

	response := map[string]any{
		"env":     app.config.env,
		"message": userWithToken,
		"status":  "true",
	}

	if err := writeJson(w, http.StatusCreated, response); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
