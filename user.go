package main

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/goadesign/goa"
	"user-microservice/app"
	"user-microservice/store"
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

func (c *UserController) Create(ctx *app.CreateUserContext) error {
    payload := goa.ContextRequest(ctx).Payload
    fmt.Println(payload)
    // Hashing
    userPassword := ctx.Payload.Password
    hashedPassword, error := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
    if error != nil {
        panic(error)
    }
    
    ctx.Payload.Password = string(hashedPassword)
    
    // Generate ID
    // ctx.Payload.ID = bson.NewObjectId()
    py := &app.Users {
        ID:   bson.NewObjectId().Hex(),
        Username: ctx.Payload.Username,
        Email: ctx.Payload.Email,
        ExternalID: ctx.Payload.ExternalID,
        Roles: ctx.Payload.Roles,
    }
    // Insert Datas
    err := c.usersc.Insert(py)
	
    // Handle errors
    if err != nil {
            if mgo.IsDup(err) {
                return ctx.BadRequest(goa.ErrBadRequest(err, "Email or Username already exists in the database"))
            }
            return err
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
