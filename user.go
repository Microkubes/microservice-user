package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/JormungandrK/backends"
	"github.com/JormungandrK/microservice-security/auth"
	"github.com/JormungandrK/microservice-user/app"
	"github.com/JormungandrK/microservice-user/store"
	"github.com/goadesign/goa"

	"golang.org/x/crypto/bcrypt"
)

// Email holds info for the email template
type Email struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
}

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	Store store.User
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, store store.User) *UserController {
	return &UserController{
		Controller: service.NewController("UserController"),
		Store:      store,
	}
}

// Create runs the create action.
func (c *UserController) Create(ctx *app.CreateUserContext) error {

	if ctx.Payload.Password == nil && ctx.Payload.ExternalID == nil {
		return ctx.BadRequest(goa.ErrBadRequest("password or externalID must be specified!"))
	}

	if len(ctx.Payload.Roles) == 0 {
		ctx.Payload.Roles = append(ctx.Payload.Roles, "user")
	}

	// Hashing password
	if ctx.Payload.Password != nil {
		hashedPassword, err := stringToBcryptHash(*ctx.Payload.Password)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		ctx.Payload.Password = &hashedPassword
	}

	result, err := c.Store.Users.Save(ctx.Payload, nil)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	if ctx.Payload.Token == nil {
		token := generateToken(42)
		ctx.Payload.Token = &token
	}

	tokenPayload := map[string]interface{}{
		"email": ctx.Payload.Email,
		"token": ctx.Payload.Token,
	}

	_, err = c.Store.Tokens.Save(&tokenPayload, nil)
	if err != nil {
		return ctx.InternalServerError(err)
	}

	user := &app.Users{}
	if err = backends.MapToInterface(result, user); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.Created(user)
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {

	user := &app.Users{}

	filter := map[string]interface{}{
		"id": ctx.UserID,
	}

	// Return one user by id.
	if err := c.Store.Users.GetOne(filter, user); err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(user)
}

// GetMe runs the getMe action.
// Get the userID from the auth ibject and return the authenticated user
func (c *UserController) GetMe(ctx *app.GetMeUserContext) error {

	var authObj *auth.Auth

	hasAuth := auth.HasAuth(ctx)

	if hasAuth {
		authObj = auth.GetAuth(ctx.Context)
	} else {
		return ctx.InternalServerError(goa.ErrInternal("Auth has not been set"))
	}

	userID := authObj.UserID

	user := &app.Users{}
	filter := map[string]interface{}{
		"id": userID,
	}

	if err := c.Store.Users.GetOne(filter, user); err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(user)
}

// Update runs the update action.
func (c *UserController) Update(ctx *app.UpdateUserContext) error {

	payload := map[string]interface{}{}

	payload["active"] = ctx.Payload.Active

	if ctx.Payload.Email != nil {
		payload["email"] = ctx.Payload.Email
	}
	if ctx.Payload.ExternalID != nil {
		payload["externalId"] = ctx.Payload.ExternalID
	}

	if ctx.Payload.Password != nil && *ctx.Payload.Password != "" {
		hashedPassword, err := stringToBcryptHash(*ctx.Payload.Password)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		payload["password"] = hashedPassword
	}

	if ctx.Payload.Roles != nil {
		payload["roles"] = ctx.Payload.Roles
	}

	if ctx.Payload.Organizations != nil {
		payload["organizations"] = ctx.Payload.Organizations
	}

	if ctx.Payload.Namespaces != nil {
		payload["namespaces"] = ctx.Payload.Namespaces
	}

	filter := map[string]interface{}{
		"id": ctx.UserID,
	}

	result, err := c.Store.Users.Save(&payload, filter)

	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	user := &app.Users{}
	if err = backends.MapToInterface(result, user); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(user)
}

// Find looks up a user by its email and password. Intended for internal use.
func (c *UserController) Find(ctx *app.FindUserContext) error {

	var userData map[string]interface{}
	filter := map[string]interface{}{
		"email": ctx.Payload.Email,
	}

	if err := c.Store.Users.GetOne(filter, &userData); err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	if _, ok := userData["password"]; !ok {
		return ctx.NotFound(goa.ErrNotFound("not found"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(ctx.Payload.Password)); err != nil {
		fmt.Println(ctx.Payload.Password)
		return ctx.NotFound(goa.ErrNotFound("not found"))
	}

	user := &app.Users{}
	if err := backends.MapToInterface(userData, user); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(user)
}

// FindByEmail looks up a user by its email.
func (c *UserController) FindByEmail(ctx *app.FindByEmailUserContext) error {

	user := &app.Users{}
	filter := map[string]interface{}{
		"email": ctx.Payload.Email,
	}

	if err := c.Store.Users.GetOne(filter, user); err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	return ctx.OK(user)
}

// Verify performs verification of the given token and activates the user for which
// this activation token was generated.
func (c *UserController) Verify(ctx *app.VerifyUserContext) error {

	if ctx.Token == nil {
		return ctx.BadRequest(goa.ErrBadRequest("token is missing from the payload"))
	}

	user := &app.Users{}
	tokenFilter := map[string]interface{}{
		"token": *ctx.Token,
	}

	err := c.Store.Tokens.GetOne(tokenFilter, user)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	emailFilter := map[string]interface{}{
		"email": user.Email,
	}
	update := map[string]interface{}{
		"active": true,
	}

	_, err = c.Store.Users.Save(&update, emailFilter)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	err = c.Store.Tokens.DeleteOne(tokenFilter)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	// empty response
	var resp []byte
	return ctx.OK(resp)
}

// ResetVerificationToken resets a verification token for a given user (by email). Generates a new value for the token
// and resets the expiration time for the token.
func (c *UserController) ResetVerificationToken(ctx *app.ResetVerificationTokenUserContext) error {

	user := &app.Users{}
	emailFilter := map[string]interface{}{
		"email": ctx.Payload.Email,
	}

	err := c.Store.Users.GetOne(emailFilter, user)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
		return ctx.InternalServerError(err)
	}

	if user == nil {
		return ctx.NotFound(fmt.Errorf("not-found"))
	}

	if user.Active {
		return ctx.BadRequest(goa.ErrBadRequest("already active"))
	}

	if err := c.Store.Tokens.DeleteOne(emailFilter); err != nil {
		if err != nil {
			e := err.(*goa.ErrorResponse)

			switch e.Status {
			case 400:
				return ctx.BadRequest(err)
			case 500:
				return ctx.InternalServerError(err)
			}
		}
	}

	token := generateToken(42)

	tokenPayload := map[string]interface{}{
		"email": ctx.Payload.Email,
		"token": token,
	}

	result, err := c.Store.Tokens.Save(&tokenPayload, nil)
	if err != nil {
		return ctx.InternalServerError(err)
	}

	resetToken := &app.ResetToken{}
	if err = backends.MapToInterface(result, resetToken); err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.OK(resetToken)
}

// generateToken generates random string with length of n
func generateToken(n int) string {
	rv := make([]byte, n)
	if _, err := rand.Reader.Read(rv); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(rv)
}

// stringToBcryptHash returns the bcrypt hash of the password with the default cost
func stringToBcryptHash(str string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedString), nil
}
