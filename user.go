package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/JormungandrK/microservice-security/auth"
	"github.com/JormungandrK/user-microservice/app"
	"github.com/JormungandrK/user-microservice/store"
	"github.com/goadesign/goa"
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
	collections store.Collections
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, collections store.Collections) *UserController {
	return &UserController{
		Controller:  service.NewController("UserController"),
		collections: collections,
	}
}

// Create runs the create action.
func (c *UserController) Create(ctx *app.CreateUserContext) error {
	if len(ctx.Payload.Roles) == 0 {
		ctx.Payload.Roles = append(ctx.Payload.Roles, "user")
	}

	id, err := c.collections.Users.CreateUser(ctx.Payload)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 400:
			return ctx.BadRequest(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	err = c.collections.Tokens.CreateToken(ctx.Payload)
	if err != nil {
		return ctx.InternalServerError(err)
	}

	var externalID string
	if ctx.Payload.ExternalID == nil {
		externalID = ""
	} else {
		externalID = *ctx.Payload.ExternalID
	}

	// Define user media type
	py := &app.Users{
		ID:            *id,
		Email:         ctx.Payload.Email,
		ExternalID:    externalID,
		Roles:         ctx.Payload.Roles,
		Organizations: ctx.Payload.Organizations,
	}

	return ctx.Created(py)
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {
	// Build the resource using the generated data structure.
	res := &app.Users{}

	// Return one user by id.
	if err := c.collections.Users.FindByID(ctx.UserID, res); err != nil {
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

	res.ID = ctx.UserID

	return ctx.OK(res)
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
	res := &app.Users{}

	if err := c.collections.Users.FindByID(userID, res); err != nil {
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

	res.ID = userID

	return ctx.OK(res)
}

// Update runs the update action.
func (c *UserController) Update(ctx *app.UpdateUserContext) error {
	res, err := c.collections.Users.UpdateUser(ctx.UserID, ctx.Payload)

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

	return ctx.OK(res)
}

// Find looks up a user by its email and password. Intended for internal use.
func (c *UserController) Find(ctx *app.FindUserContext) error {
	user, err := c.collections.Users.FindByEmailAndPassword(ctx.Payload.Email, ctx.Payload.Password)

	if err != nil {
		fmt.Printf("Failed to find user. Error: %s", err)
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if user == nil {
		return ctx.NotFound()
	}

	return ctx.OK(user)
}

// FindByEmail looks up a user by its email.
func (c *UserController) FindByEmail(ctx *app.FindByEmailUserContext) error {
	user, err := c.collections.Users.FindByEmail(ctx.Payload.Email)
	if err != nil {
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

func (c *UserController) Verify(ctx *app.VerifyUserContext) error {
	// Check if user is already activated
	token := *ctx.Token
	email, err := c.collections.Tokens.VerifyToken(token)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	err = c.collections.Users.ActivateUser(*email)
	if err != nil {
		e := err.(*goa.ErrorResponse)

		switch e.Status {
		case 404:
			return ctx.NotFound(err)
		default:
			return ctx.InternalServerError(err)
		}
	}

	err = c.collections.Tokens.DeleteToken(token)
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

func (c *UserController) ResetVerificationToken(ctx *app.ResetVerificationTokenUserContext) error {
	user, err := c.collections.Users.FindByEmail(ctx.Payload.Email)
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.NotFound(err)
		}
		return ctx.InternalServerError(err)
	}

	if user.Active {
		return ctx.BadRequest(goa.ErrBadRequest("already active"))
	}

	if err := c.collections.Tokens.DeleteUserToken(user.ID); err != nil {
		return ctx.InternalServerError(err)
	}

	token := generateToken(42)

	if err := c.collections.Tokens.CreateToken(&app.UserPayload{
		Active:        false,
		Email:         user.Email,
		ExternalID:    &user.ExternalID,
		Organizations: user.Organizations,
		Roles:         user.Roles,
		Token:         &token,
	}); err != nil {
		return ctx.InternalServerError(err)
	}

	return ctx.OK(&app.ResetToken{
		Email: user.Email,
		ID:    user.ID,
		Token: token,
	})
}

func generateToken(n int) string {
	rv := make([]byte, n)
	if _, err := rand.Reader.Read(rv); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(rv)
}
