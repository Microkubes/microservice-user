package main

import (
	"net/http"
	"os"

	"github.com/Microkubes/backends"
	"github.com/Microkubes/microservice-security/chain"
	"github.com/Microkubes/microservice-security/flow"
	"github.com/Microkubes/microservice-tools/config"
	"github.com/Microkubes/microservice-tools/gateway"
	"github.com/Microkubes/microservice-user/app"
	"github.com/Microkubes/microservice-user/store"

	"github.com/Microkubes/microservice-tools/utils"
	"github.com/Microkubes/microservice-tools/utils/version"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("user")

	gatewayAdminURL, configFile := loadGatewaySettings()

	serviceConfig, err := config.LoadConfig(configFile)
	if err != nil {
		service.LogError("config", "err", err)
		return
	}

	registration := gateway.NewKongGateway(gatewayAdminURL, &http.Client{}, serviceConfig.Service)

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

	service.Use(healthcheck.NewCheckMiddleware("/healthcheck"))

	service.Use(version.NewVersionMiddleware(serviceConfig.Version, "/version"))

	// Get the db collections/tables
	dbConf := serviceConfig.DBConfig

	backendManager := backends.NewBackendSupport(map[string]*config.DBInfo{
		"mongodb":  &dbConf.DBInfo,
		"dynamodb": &dbConf.DBInfo,
	})

	backend, err := backendManager.GetBackend(dbConf.DBName)
	if err != nil {
		service.LogError("Failed to configure backend. ", err)
	}

	userRepo, err := backend.DefineRepository("users", backends.RepositoryDefinitionMap{
		"name": "users",
		"indexes": []backends.Index{
			backends.NewUniqueIndex("email"),
		},
		"hashKey":       "email",
		"readCapacity":  int64(5),
		"writeCapacity": int64(5),
		"GSI": map[string]interface{}{
			"email": map[string]interface{}{
				"readCapacity":  1,
				"writeCapacity": 1,
			},
		},
	})
	if err != nil {
		service.LogError("Failed to get users repo.", err)
		return
	}

	tokenRepo, err := backend.DefineRepository("tokens", backends.RepositoryDefinitionMap{
		"name": "tokens",
		"indexes": []backends.Index{
			backends.NewUniqueIndex("token"),
		},
		"hashKey":       "token",
		"readCapacity":  int64(5),
		"writeCapacity": int64(5),
		"GSI": map[string]interface{}{
			"token": map[string]interface{}{
				"readCapacity":  1,
				"writeCapacity": 1,
			},
		},
		"enableTtl":    true,
		"ttlAttribute": "created_at",
		"ttl":          86400,
	})
	if err != nil {
		service.LogError("Failed to get tokens repo.", err)
		return
	}

	store := store.User{
		Users:  userRepo,
		Tokens: tokenRepo,
	}

	// Mount "swagger" controller
	c1 := NewSwaggerController(service)
	app.MountSwaggerController(service, c1)
	// Mount "user" controller
	c2 := NewUserController(service, store)
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
