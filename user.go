package main

import (
	"fmt"
	"github.com/goadesign/goa"
	"user-microservice/app"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Person struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string
	Phone     string
	Timestamp time.Time
}

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
	usersc *mgo.Collection
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service, usersc *mgo.Collection) *UserController {
	return &UserController{
		Controller: service.NewController("UserController"),
		usersc: usersc,
	}
}

// Create runs the create action.
func (c *UserController) Create(ctx *app.CreateUserContext) error {
	// Insert Datas
	err := c.usersc.Insert(&Person{Name: "Vlado", Phone: "+55 53 1234 4321", Timestamp: time.Now()},
		&Person{Name: "Ace", Phone: "+66 33 1234 5678", Timestamp: time.Now()})

	if err != nil {
		panic(err)
	}

	return nil
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {
	// UserController_Get: start_implement
	if ctx.UserID == 0 {
		// Emulate a missing record with ID 0
		return ctx.NotFound()
	}

	roles := []string{"admin", "owner"}

	// Build the resource using the generated data structure
	user := &app.Users {
		ID:   ctx.UserID,
		Username: fmt.Sprintf("User #%d", ctx.UserID),
		Email: "example@gmail.com",
		ExternalID: "qwe23adsa213saqqw",
		Roles: roles,
	}

	// UserController_Get: end_implement
	return ctx.OK(user)
}

// GetMe runs the getMe action.
func (c *UserController) GetMe(ctx *app.GetMeUserContext) error {
	// UserController_GetMe: start_implement

	// Put your logic here

	// UserController_GetMe: end_implement
	res := &app.Users{}
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
