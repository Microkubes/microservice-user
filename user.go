package main

import (
	"time"
	"user-microservice/app"
	"user-microservice/store"
	"github.com/goadesign/goa"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"golang.org/x/crypto/bcrypt"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	usersCollection store.Collection
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, usersCollection store.Collection) *UserController {
	return &UserController{
		Controller: service.NewController("UserController"),
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
		"_id": id, 
		"username": ctx.Payload.Username,
		"email": ctx.Payload.Email,
		"password": string(hashedPassword),
		"externalId": ctx.Payload.ExternalID,
		"roles": ctx.Payload.Roles,
	})
	
	// Handle errors
	if err != nil {
		if mgo.IsDup(err) {
			return ctx.BadRequest(goa.ErrBadRequest(err, "Email or Username already exists in the database"))
		}
		return err
	}
	
	// Define user media type
	py := &app.Users {
		ID:			id.Hex(),
		Username:	ctx.Payload.Username,
		Email:		ctx.Payload.Email,
		ExternalID:	ctx.Payload.ExternalID,
		Roles:		ctx.Payload.Roles,
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
	userId := bson.ObjectIdHex(ctx.UserID)

	// Return true if userId is valid. A valid userId must contain exactly 12 bytes.
	if userId.Valid() != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid Id"))
	}

	// Return one user by id.
	if err := c.usersCollection.FindByID(userId, res); err != nil {
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
	userId := bson.ObjectIdHex(*ctx.UserID)

	// Return true if userId is valid. A valid userId must contain exactly 12 bytes.
	if userId.Valid() != true {
		return ctx.NotFound(goa.ErrNotFound("Invalid Id"))
	}

	// Return one user by id.
	if err := c.usersCollection.FindByID(userId, res); err != nil {
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
