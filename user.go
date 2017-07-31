package main

import (
	"github.com/JormungandrK/user-microservice/app"
	"github.com/JormungandrK/user-microservice/store"

	"github.com/goadesign/goa"
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
	// UserController_Create: start_implement

	// Put your logic here

	// UserController_Create: end_implement
	return nil
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

	// Return true if userId is valid. A valid userId must contain exactly 12 bytes.
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

	// Return true if userId is valid. A valid userId must contain exactly 12 bytes.
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
	// UserController_Update: start_implement

	// Put your logic here

	// UserController_Update: end_implement
	res := &app.Users{}
	return ctx.OK(res)
}
