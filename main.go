package main

import (
	"net/http"
	"os"

	"github.com/JormungandrK/microservice-security/chain"
	"github.com/JormungandrK/microservice-security/flow"
	"github.com/JormungandrK/microservice-tools/config"
	"github.com/JormungandrK/user-microservice/app"

	"github.com/JormungandrK/microservice-tools/gateway"

	"github.com/JormungandrK/user-microservice/store"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("user")

	_, configFile := loadGatewaySettings()

	serviceConfig, err := config.LoadConfig(configFile)
	if err != nil {
		service.LogError("config", "err", err)
		return
	}

	//registration, err := gateway.NewKongGatewayFromConfigFile(gatewayURL, &http.Client{}, configFile)
	registration := gateway.NewKongGateway(serviceConfig.GatewayAdminURL, &http.Client{}, serviceConfig.Service)

	err = registration.SelfRegister()
	if err != nil {
		panic(err)
	}

	defer registration.Unregister()

	securityChain, cleanup, err := flow.NewSecurityFromConfig(serviceConfig)
	if err != nil {
		panic(err)
	}

	defer cleanup()

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	service.Use(chain.AsGoaMiddleware(securityChain))

	dbConf := serviceConfig.DBConfig
	// Create new session to MongoDB
	session := store.NewSession(dbConf.Host, dbConf.Username, dbConf.Password, dbConf.DatabaseName)

	// At the end close session
	defer session.Close()

	// Create users collection and indexes
	indexes := []string{"email"}
	usersCollection := store.PrepareDB(session, dbConf.DatabaseName, "users", indexes)

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

func loadMongnoSettings() (string, string, string, string) {
	host := os.Getenv("MONGO_URL")
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
		serviceConfigFile = "/run/secrets/microservice_user_config.json"
	}

	return gatewayURL, serviceConfigFile
}
