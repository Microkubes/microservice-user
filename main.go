//go:generate goagen bootstrap -d user-microservice/design

package main

import (
	"net/http"
	"os"
	"user-microservice/app"

	"github.com/JormungandrK/microservice-tools/gateway"

	"user-microservice/store"

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

	// Load MongoDB ENV variables
	host, username, password, database := loadMongnoSettings()
	// Create new session to MongoDB
	session := store.NewSession(host, username, password, database)

	// At the end close session
	defer session.Close()

	// Create users collection and indexes
	indexes :=  []string{"username", "email"}
	usersCollection := store.PrepareDB(session, database, "users", indexes)

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
