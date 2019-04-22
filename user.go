package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Microkubes/backends"
	"github.com/Microkubes/microservice-security/auth"
	"github.com/Microkubes/microservice-tools/rabbitmq"
	"github.com/Microkubes/microservice-user/app"
	"github.com/Microkubes/microservice-user/store"
	"github.com/goadesign/goa"

	"golang.org/x/crypto/bcrypt"
)

// Email holds info for the email template
type Email struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
}

// EmailInfo holds data for "verification-email" RabbitMQ channel
type EmailInfo struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Template string `json:"template,omitempty"`
	Email    string `json:"email,omitempty"`
	Token    string `json:"token,omitempty"`
}

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	Store           store.User
	ChannelRabbitMQ rabbitmq.Channel
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, store store.User, rmqChannel rabbitmq.Channel) *UserController {
	return &UserController{
		Controller:      service.NewController("UserController"),
		Store:           store,
		ChannelRabbitMQ: rmqChannel,
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

	user := &store.UserRecord{
		Active: false,
		Email:  ctx.Payload.Email,
		//ExternalID:    ctx.Payload.ExternalID == nil ? "": *ctx.Payload.ExternalID,
		Namespaces:    ctx.Payload.Namespaces,
		Organizations: ctx.Payload.Organizations,
		Password:      *ctx.Payload.Password,
		Roles:         ctx.Payload.Roles,
		//Token:         ctx.Payload.Token,
	}

	if ctx.Payload.ExternalID != nil {
		user.ExternalID = *ctx.Payload.ExternalID
	}
	if ctx.Payload.Token != nil {
		user.Token = *ctx.Payload.Token
	}

	result, err := c.Store.Users.Save(user, nil)
	if err != nil {
		if backends.IsErrAlreadyExists(err) || backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
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

	return ctx.Created(result.(*store.UserRecord).ToAppUsers())
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {

	user := &app.Users{}

	// Return one user by id.
	if _, err := c.Store.Users.GetOne(backends.NewFilter().Match("id", ctx.UserID), user); err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(user)
}

// GetMe runs the getMe action.
// Get the userID from the auth ibject and return the authenticated user
func (c *UserController) GetMe(ctx *app.GetMeUserContext) error {

	if !auth.HasAuth(ctx.Context) {
		return ctx.InternalServerError(goa.ErrBadRequest("no-auth"))
	}

	userID := auth.GetAuth(ctx.Context).UserID

	user := &app.Users{}
	if _, err := c.Store.Users.GetOne(backends.NewFilter().Match("id", userID), user); err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(user)
}

//GetAll retrives all active users
func (c *UserController) GetAll(ctx *app.GetAllUserContext) error {
	if !auth.HasAuth(ctx.Context) {
		return ctx.InternalServerError(goa.ErrBadRequest("no-auth"))
	}

	order := ""
	sorting := ""
	limit := 0
	offset := 0

	if ctx.Order != nil {
		order = *ctx.Order
	}

	if ctx.Offset != nil {
		offset = *ctx.Offset
	}

	if ctx.Limit != nil {
		limit = *ctx.Limit
	}

	if ctx.Sorting != nil {
		sorting = *ctx.Sorting
	}

	var typeHint map[string]interface{}
	users, err := c.Store.Users.GetAll(backends.NewFilter().Match("active", true), typeHint, order, sorting, limit, offset)
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err.Error()))
		}

		return ctx.InternalServerError(goa.ErrInternal(err.Error()))
	}

	resp, err := json.Marshal(users)
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err.Error()))
	}

	return ctx.OK(resp)
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

	result, err := c.Store.Users.Save(&payload, backends.NewFilter().Match("id", ctx.UserID))

	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
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
	if _, err := c.Store.Users.GetOne(backends.NewFilter().Match("email", ctx.Payload.Email), &userData); err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}

		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
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
	if _, err := c.Store.Users.GetOne(backends.NewFilter().Match("email", ctx.Payload.Email), user); err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
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
	_, err := c.Store.Tokens.GetOne(backends.NewFilter().Match("token", *ctx.Token), user)
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	update := map[string]interface{}{
		"active": true,
	}
	_, err = c.Store.Users.Save(&update, backends.NewFilter().Match("email", user.Email))
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	err = c.Store.Tokens.DeleteOne(backends.NewFilter().Match("token", *ctx.Token))
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	// empty response
	var resp []byte
	return ctx.OK(resp)
}

// ResetVerificationToken resets a verification token for a given user (by email). Generates a new value for the token
// and resets the expiration time for the token.
func (c *UserController) ResetVerificationToken(ctx *app.ResetVerificationTokenUserContext) error {

	user := &app.Users{}
	_, err := c.Store.Users.GetOne(backends.NewFilter().Match("email", ctx.Payload.Email), user)
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.NotFound(goa.ErrNotFound(err))
		}
		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if user == nil {
		return ctx.NotFound(fmt.Errorf("not-found"))
	}

	if user.Active {
		return ctx.BadRequest(goa.ErrBadRequest("already active"))
	}

	if err := c.Store.Tokens.DeleteOne(backends.NewFilter().Match("email", ctx.Payload.Email)); err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
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

// ForgotPassword is used for verifying user and sending mail with generated token
func (c *UserController) ForgotPassword(ctx *app.ForgotPasswordUserContext) error {
	userRecord := &store.UserRecord{}
	_, err := c.Store.Users.GetOne(backends.NewFilter().Match("email", ctx.Payload.Email), userRecord)
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.OK([]byte{})
		}
		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}
	fpToken := store.FPToken{}
	fpToken.Token = generateToken(42)
	fpToken.ExpDate = generateExpDate()
	userRecord.FPToken = fpToken
	_, err = c.Store.Users.Save(userRecord, backends.NewFilter().Match("id", userRecord.ID))
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	fmt.Println("sending mail... TOKEN: " + fpToken.Token)

	if c.ChannelRabbitMQ != nil {

		emailInfo := EmailInfo{
			ID:       userRecord.ID,
			Name:     "User",
			Email:    userRecord.Email,
			Token:    fpToken.Token,
			Template: "Forgot Password",
		}

		body, err := json.Marshal(emailInfo)
		if err != nil {
			c.Service.LogError("User: failed to serialize emailInfo.", "err", err.Error())
		}

		if err := c.ChannelRabbitMQ.Send("verification-email", body); err != nil {
			c.Service.LogError("User: failed to send message on rabbitMQ.", "err", err.Error())
			return ctx.InternalServerError(goa.ErrInternal(err))
		}
	}

	return ctx.OK([]byte{})
}

// ForgotPasswordUpdate endpoint for changing old password with new one
func (c *UserController) ForgotPasswordUpdate(ctx *app.ForgotPasswordUpdateUserContext) error {
	userRecord := &store.UserRecord{}
	_, err := c.Store.Users.GetOne(backends.NewFilter().Match("email", ctx.Payload.Email), userRecord)
	if err != nil {
		if backends.IsErrNotFound(err) {
			return ctx.OK([]byte{})
		}
		if backends.IsErrInvalidInput(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err))
		}
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	if userRecord.FPToken.Token != ctx.Payload.Token {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}
	if !checkExpDate(userRecord.FPToken.ExpDate) {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	hashedPassword, err := stringToBcryptHash(ctx.Payload.Password)
	if err != nil {
		return ctx.BadRequest(goa.ErrBadRequest(err))
	}

	userRecord.FPToken.ExpDate = "0"
	userRecord.Password = hashedPassword

	_, err = c.Store.Users.Save(userRecord, backends.NewFilter().Match("id", userRecord.ID))
	if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte{})
}

// generateToken generates random string with length of n
func generateToken(n int) string {
	rv := make([]byte, n)
	if _, err := rand.Reader.Read(rv); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(rv)
}

func generateExpDate() string {
	validTime := 60 * 24 // 1 day valid time
	time := int(time.Now().UTC().Unix()/60) + validTime
	return strconv.Itoa(time)
}

func checkExpDate(expDate string) bool {
	fmt.Println(expDate)
	expDateParsed, err := strconv.Atoi(expDate)
	if err != nil {
		return false
	}
	currentTime := int(time.Now().UTC().Unix() / 60)
	if expDateParsed > currentTime {
		return true
	}
	return false
}

// stringToBcryptHash returns the bcrypt hash of the password with the default cost
func stringToBcryptHash(str string) (string, error) {
	hashedString, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedString), nil
}
