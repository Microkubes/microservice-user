package main

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	// "encoding/json"

	"github.com/JormungandrK/backends"
	"github.com/JormungandrK/microservice-security/auth"
	"github.com/JormungandrK/user-microservice/app"
	// "github.com/JormungandrK/user-microservice/store"
	"github.com/goadesign/goa"

	"golang.org/x/crypto/bcrypt"
)

// Email holds info for the email template
type Email struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
}

// UserStore wraps User's collections/tables. Implements backneds.Repository interface
type UserStore struct {
	Users  backends.Repository
	Tokens backends.Repository
}

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	Store UserStore
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, store UserStore) *UserController {
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

	if ctx.Payload.Token == nil {
		token := generateToken(42)
		ctx.Payload.Token = &token
	}

	// Hashing password
	if ctx.Payload.Password != nil {
		userPassword := *ctx.Payload.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
		pass := string(hashedPassword)
		ctx.Payload.Password = &pass
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

	tokenPayload := map[string]interface{}{
		"email":      ctx.Payload.Email,
		"token":      ctx.Payload.Token,
		"created_at": time.Now(),
	}

	_, err = c.Store.Tokens.Save(&tokenPayload, nil)
	if err != nil {
		return ctx.InternalServerError(err)
	}

	var externalID string
	if ctx.Payload.ExternalID == nil {
		externalID = ""
	} else {
		externalID = *ctx.Payload.ExternalID
	}

	r := result.(map[string]interface{})
	py := &app.Users{
		ID:            r["id"].(string),
		Email:         ctx.Payload.Email,
		ExternalID:    externalID,
		Roles:         ctx.Payload.Roles,
		Organizations: ctx.Payload.Organizations,
	}

	return ctx.Created(py)
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

	payload := map[string]interface{}{
		"id": ctx.UserID,
	}

	payload["active"] = ctx.Payload.Active
	if ctx.Payload.Email != "" {
		payload["email"] = ctx.Payload.Email
	}
	if ctx.Payload.ExternalID != nil {
		payload["externalId"] = ctx.Payload.ExternalID
	}

	if ctx.Payload.Password != nil && *ctx.Payload.Password != "" {
		hashedPassword, herr := bcrypt.GenerateFromPassword([]byte(*ctx.Payload.Password), bcrypt.DefaultCost)
		if herr != nil {
			return goa.ErrInternal(herr)
		}
		payload["password"] = string(hashedPassword)
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
	err = backends.MapToInterface(result, user)
	if err != nil {
		return ctx.InternalServerError(err)
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
		return ctx.NotFound(goa.ErrNotFound("not found"))
	}

	user := &app.Users{}
	err := backends.MapToInterface(userData, user)
	if err != nil {
		return ctx.InternalServerError(err)
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
	// Check if user is already activated
	// token := *ctx.Token
	// email, err := c.collections.Tokens.VerifyToken(token)
	// if err != nil {
	// 	e := err.(*goa.ErrorResponse)

	// 	switch e.Status {
	// 	case 404:
	// 		return ctx.NotFound(err)
	// 	default:
	// 		return ctx.InternalServerError(err)
	// 	}
	// }

	// err = c.collections.Users.ActivateUser(*email)
	// if err != nil {
	// 	e := err.(*goa.ErrorResponse)

	// 	switch e.Status {
	// 	case 404:
	// 		return ctx.NotFound(err)
	// 	default:
	// 		return ctx.InternalServerError(err)
	// 	}
	// }

	// err = c.collections.Tokens.DeleteToken(token)
	// if err != nil {
	// 	e := err.(*goa.ErrorResponse)

	// 	switch e.Status {
	// 	case 404:
	// 		return ctx.NotFound(err)
	// 	default:
	// 		return ctx.InternalServerError(err)
	// 	}
	// }

	// // empty response
	// var resp []byte
	// return ctx.OK(resp)

	return nil
}

// ResetVerificationToken resets a verification token for a given user (by email). Generates a new value for the token
// and resets the expiration time for the token.
func (c *UserController) ResetVerificationToken(ctx *app.ResetVerificationTokenUserContext) error {
	// user, err := c.collections.Users.FindByEmail(ctx.Payload.Email)
	// if err != nil {
	// 	if err.Error() == "user not found" {
	// 		return ctx.NotFound(err)
	// 	}
	// 	return ctx.InternalServerError(err)
	// }
	// if user == nil {
	// 	return ctx.NotFound(fmt.Errorf("not-found"))
	// }

	// if user.Active {
	// 	return ctx.BadRequest(goa.ErrBadRequest("already active"))
	// }

	// if user.ExternalID != "" {
	// 	return ctx.BadRequest(goa.ErrBadRequest("external-user"))
	// }

	// if err := c.collections.Tokens.DeleteUserToken(user.Email); err != nil {
	// 	return ctx.InternalServerError(err)
	// }

	// token := generateToken(42)

	// if err := c.collections.Tokens.CreateToken(&app.UserPayload{
	// 	Active:        false,
	// 	Email:         user.Email,
	// 	ExternalID:    &user.ExternalID,
	// 	Organizations: user.Organizations,
	// 	Roles:         user.Roles,
	// 	Token:         &token,
	// }); err != nil {
	// 	return ctx.InternalServerError(err)
	// }

	// return ctx.OK(&app.ResetToken{
	// 	Email: user.Email,
	// 	ID:    user.ID,
	// 	Token: token,
	// })

	return nil
}

func generateToken(n int) string {
	rv := make([]byte, n)
	if _, err := rand.Reader.Read(rv); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(rv)
}
