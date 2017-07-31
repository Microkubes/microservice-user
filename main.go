//go:generate goagen bootstrap -d user-microservice/design

package main

import (
	"net/http"
	"os"

	"github.com/JormungandrK/user-microservice/app"

	"github.com/JormungandrK/microservice-tools/gateway"

	"github.com/JormungandrK/user-microservice/store"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
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

	// Create new session to MongoDB
	session := store.NewSession()

	// At the end close session
	defer session.Close()

	// Create users collection and indexes
	indexes := []string{"username", "email"}
	usersCollection := store.PrepareDB(session, "users", "users", indexes)

	// Mount "swagger" controller
	c1 := NewSwaggerController(service)
	app.MountSwaggerController(service, c1)
	// Mount "user" controller
	c2 := NewUserController(service, &store.MongoCollection{Collection: usersCollection})
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
