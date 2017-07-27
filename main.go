//go:generate goagen bootstrap -d user-microservice/design

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"user-microservice/store"
	"user-microservice/app"
)

func main() {
	// Create service
	service := goa.New("user")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Create new session to MongoDB
	session := store.NewSession()

	// At the end close session
	defer session.Close()

	// Create users collection and indexes
	indexes :=  []string{"username", "email"}
	usersCollection := store.PrepareDB(session, "users", "users", indexes)

	// Mount "swagger" controller
	c1 := NewSwaggerController(service)
	app.MountSwaggerController(service, c1)
	// Mount "user" controller
	c2 := NewUserController(service, &store.MongoCollection{usersCollection})
	app.MountUserController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}
