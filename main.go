//go:generate goagen bootstrap -d user-microservice/design

package main

import (
	"microservice-tools/gateway"
	"net/http"
	"os"
	"user-microservice/app"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"gopkg.in/mgo.v2"
)

const (
	Host     = "127.0.0.1:27017"
	Username = "restapi"
	Password = "restapi"
	Database = "users"
)

func main() {
	gatewayURL, configFile := loadGatewaySettings()
	registration, err := gateway.NewKongGatewayFromConfigFile(gatewayURL, &http.Client{}, configFile)
	if err != nil {
		panic(err)
	}
	err = registration.SelfRegister()
	if err != nil {
		panic(err)
	}

	defer registration.Unregister()

	// Create service
	service := goa.New("user")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Connnect to MongoDB
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{Host},
		Username: Username,
		Password: Password,
		Database: Database,
	})
	if err != nil {
		panic(err)
	}

	// At the end close session
	defer session.Close()

	// SetMode - consistency mode for the session.
	session.SetMode(mgo.Monotonic, true)

	// Create usersCollection collection
	usersCollection := session.DB("users").C("users")

	// Define indexes
	index := mgo.Index{
		Key:        []string{"username", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	// Create indexes
	if err = usersCollection.EnsureIndex(index); err != nil {
		panic(err)
	}

	// Mount "swagger" controller
	c1 := NewSwaggerController(service)
	app.MountSwaggerController(service, c1)
	// Mount "user" controller
	c2 := NewUserController(service, usersCollection)
	app.MountUserController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}

}

func loadGatewaySettings() (string, string) {
	gatewayURL := os.Getenv("API_GATEWAY_URL")
	serviceConfigFile := os.Getenv("SERVICE_CONFIG_FILE")

	if gatewayURL == "" {
		gatewayURL = "http://localhost:8001"
	}
	if serviceConfigFile == "" {
		serviceConfigFile = "config.json"
	}

	return gatewayURL, serviceConfigFile
}
