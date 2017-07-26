package main

import (
	"fmt"
	"user-microservice/app"

	"github.com/goadesign/goa"
)

// UserController implements the user resource.
type UserController struct {
	*goa.Controller
}

// NewUserController creates a user controller.
func NewUserController(service *goa.Service) *UserController {
	return &UserController{Controller: service.NewController("UserController")}
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
	// UserController_Get: start_implement
	if ctx.UserID == 0 {
		// Emulate a missing record with ID 0
		return ctx.NotFound()
	}

	roles := []string{"admin", "owner"}

	// Build the resource using the generated data structure
	user := &app.Users{
		ID:         ctx.UserID,
		Username:   fmt.Sprintf("User #%d", ctx.UserID),
		Email:      "example@gmail.com",
		ExternalID: "qwe23adsa213saqqw",
		Roles:      roles,
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
