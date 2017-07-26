package main

import (
	"github.com/goadesign/goa"
	"user-microservice/app"
	"user-microservice/store"
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
	// UserController_Create: start_implement

	// Put your logic here

	// UserController_Create: end_implement

	return nil
}

// Get runs the get action.
func (c *UserController) Get(ctx *app.GetUserContext) error {
	// Build the resource using the generated data structure.
	res := &app.Users{}

	// Return one user by id.
	if err := c.usersCollection.FindByID(ctx.UserID, res); err != nil {
		return ctx.NotFound(goa.ErrNotFound(err))
	}

	res.ID = ctx.UserID
	
	return ctx.OK(res)
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
