package main

import (
	"github.com/JormungandrK/user-microservice/app"
	"github.com/JormungandrK/user-microservice/store"

	"time"

	"github.com/goadesign/goa"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	usersCollection store.Collection
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, usersCollection store.Collection) *UserController {
	return &UserController{
		Controller:      service.NewController("UserController"),
		usersCollection: usersCollection,
	}
}

// Create runs the create action.
func (c *UserController) Create(ctx *app.CreateUserContext) error {
	// Hashing
	userPassword := ctx.Payload.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert Data
	id := bson.NewObjectIdWithTime(time.Now())
	err = c.usersCollection.Insert(bson.M{
		"_id":        id,
		"username":   ctx.Payload.Username,
		"email":      ctx.Payload.Email,
		"password":   string(hashedPassword),
		"externalId": ctx.Payload.ExternalID,
		"roles":      ctx.Payload.Roles,
	})

	// Handle errors
	if err != nil {
		if mgo.IsDup(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err, "Email or Username already exists in the database"))
		}
		return err
	}

	// Define user media type
	py := &app.Users{
		ID:         id.Hex(),
		Username:   ctx.Payload.Username,
		Email:      ctx.Payload.Email,
		ExternalID: ctx.Payload.ExternalID,
		Roles:      ctx.Payload.Roles,
	}

	return ctx.Created(py)
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {
	// Build the resource using the generated data structure.
	res := &app.Users{}

	// Return whether ctx.UserID is a valid hex representation of an ObjectId.
	if bson.IsObjectIdHex(ctx.UserID) != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid Id"))
	}

	// Return an ObjectId from the provided hex representation.
	userID := bson.ObjectIdHex(ctx.UserID)

	// Return true if userID is valid. A valid userID must contain exactly 12 bytes.
	if userID.Valid() != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid Id"))
	}

	// Return one user by id.
	if err := c.usersCollection.FindByID(userID, res); err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	res.ID = ctx.UserID

	return ctx.OK(res)
}

// GetMe runs the getMe action.
func (c *UserController) GetMe(ctx *app.GetMeUserContext) error {
	// Build the resource using the generated data structure.
	res := &app.Users{}

	// Return whether ctx.UserID is a valid hex representation of an ObjectId.
	if bson.IsObjectIdHex(*ctx.UserID) != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid Id"))
	}

	// Return an ObjectId from the provided hex representation.
	userID := bson.ObjectIdHex(*ctx.UserID)

	// Return true if userID is valid. A valid userID must contain exactly 12 bytes.
	if userID.Valid() != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid Id"))
	}

	// Return one user by id.
	if err := c.usersCollection.FindByID(userID, res); err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	res.ID = *ctx.UserID

	return ctx.OK(res)
}

// Update runs the update action.
func (c *UserController) Update(ctx *app.UpdateUserContext) error {
	email := ctx.Payload.Email
	password := ctx.Payload.Password
	roles := ctx.Payload.Roles
	username := ctx.Payload.Username
	id := ctx.UserID

	// Update
	docId := bson.M{"_id": bson.ObjectIdHex(id)}
	change := bson.M{"$set": bson.M{"username": username, "roles": roles, "password": password, "email": email}}

	err := c.usersCollection.Update(docId, change)
	if err != nil {
		return err
	}

	res := &app.Users{}

	if err = c.usersCollection.FindByID(bson.ObjectIdHex(id), res); err != nil {
		return ctx.Err()
	}
	return ctx.OK(res)
}
