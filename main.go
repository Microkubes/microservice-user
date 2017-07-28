//go:generate goagen bootstrap -d user-microservice/design

package main

import (
	"os"
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

	// Load MongoDB ENV variables
	host, username, password, database := loadMongnoSettings()
	// Create new session to MongoDB
	session := store.NewSession(host, usersname, password, database)

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

func loadMongnoSettings() (string, string, string, string) {
	host     := os.Getenv("MONGO_URL")
    username := os.Getenv("MS_USERNAME")
    password := os.Getenv("MS_PASSWORD")
    database := os.Getenv("MS_DBNAME")

    if host == "" {
    	host = "127.0.0.1:27017"
    }
    if username == "" {
    	username = "restapi"
    }
    if password == "" {
    	password = "restapi"
    }
    if database == "" {
    	database = "users"
    }

    return host, username, password, database
}