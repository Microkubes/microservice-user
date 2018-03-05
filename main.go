package main

import (
	"net/http"
	"os"

	"github.com/JormungandrK/backends"
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

	// backends

	backendManager := backends.NewBackendManager(map[string]*config.DBInfo{})

	backend, err := backendManager.GetBackend("mongodb") // from config
	if err != nil {
		service.LogError("Failed to configure backend. ", err)
	}

	userRepo, err := backend.DefineRepository("user", backends.RepositoryDefinitionMap{
		"indexes": []string{"id", "email"},
	})
	if err != nil {
		service.LogError("Failed to get users repo.", err)
	}

	tokenRepo, err := backend.DefineRepository("token", backends.RepositoryDefinitionMap{
		"indexes":   []string{"id", "token", "etc"},
		"enableTtl": true,
		"ttl":       86400,
	})

	if err != nil {
		service.LogError("Failed to get tokens repo.", err)
	}

	service.LogInfo("Set up users and tokens repositories:", userRepo, tokenRepo)

	// Later on, we can obtain the repositories like so:

	repo, err := backend.GetRepository("user")
	if err != nil {
		// the repo is not defined
		service.LogError("Repo is not defined", err)
	}
	service.LogInfo("Got user repo: ", repo)

	// end backends conf

	dbConf := serviceConfig.DBConfig
	// Create new session to MongoDB
	session := store.NewSession(dbConf.Host, dbConf.Username, dbConf.Password, dbConf.DatabaseName)

	// At the end close session
	defer session.Close()

	// Create users collection and indexes
	indexes := []string{"email"}
	usersCollection := store.PrepareDB(session, dbConf.DatabaseName, "users", indexes, false)
	indexes = []string{"token"}
	tokensCollection := store.PrepareDB(session, dbConf.DatabaseName, "tokens", indexes, true)

	// Mount "swagger" controller
	c1 := NewSwaggerController(service)
	app.MountSwaggerController(service, c1)
	// Mount "user" controller
	c2 := NewUserController(service, store.Collections{Users: &store.UserCollection{Collection: usersCollection}, Tokens: &store.TokenCollection{Collection: tokensCollection}})
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
