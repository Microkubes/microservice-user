package design

// Use . imports to enable the DSL
import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

// Define default description and default global property values
var _ = API("user", func() {
	Title("The user microservice")
	Description("A service that provides basic access to the user data")
	Version("1.0")
	Scheme("http")
	Host("localhost:8080")
})

// Resources group related API endpoints together.
var _ = Resource("user", func() {
	BasePath("/users")
	DefaultMedia(UserMedia)
	// Do not setup security here!

	// Actions define a single API endpoint
	Action("create", func() {
		Description("Creates user")
		Routing(POST(""))
		Payload(UserPayload)
		Response(Created, UserMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("get", func() {
		Description("Get user by id")
		Routing(GET("/:userId"))
		Params(func() {
			Param("userId", String, "User ID")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getMe", func() {
		Description("Retrieves the user information for the authenticated user")
		Routing(GET("/me"))
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("update", func() {
		Description("Update user")
		Routing(PUT("/:userId"))
		Params(func() {
			Param("userId", String, "User ID")
		})
		Payload(UserPayload)
		Response(OK, UserMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("find", func() {
		Description("Find a user by username+password")
		Routing(POST("find"))
		Payload(CredentialsPayload)
		Response(OK, UserMedia)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("findByEmail", func() {
		Description("Find a user by email")
		Routing(POST("find/email"))
		Payload(EmailPayload)
		Response(OK, UserMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})

// UserMedia defines the media type used to render user.
var UserMedia = MediaType("application/vnd.goa.user+json", func() {
	TypeName("users")
	Reference(UserPayload)

	Attributes(func() {
		Attribute("id", String, "Unique user ID")
		Attribute("username")
		Attribute("email")
		Attribute("roles")
		Attribute("externalId")
		Attribute("active")
		Required("id", "username", "email", "roles", "externalId", "active")
	})

	View("default", func() {
		Attribute("id")
		Attribute("username")
		Attribute("email")
		Attribute("roles")
		Attribute("externalId")
		Attribute("active")
	})
})

// UserPayload defines the payload for the user.
var UserPayload = Type("UserPayload", func() {
	Description("UserPayload")

	Attribute("username", String, "Name of user", func() {
		MinLength(4)
		MaxLength(50)
	})
	Attribute("email", String, "Email of user", func() {
		Format("email")
	})
	Attribute("password", String, "Password of user", func() {
		MinLength(6)
		MaxLength(30)
	})
	Attribute("roles", ArrayOf(String), "Roles of user")
	Attribute("externalId", String, "External id of user")
	Attribute("active", Boolean, "Status of user account", func() {
		Default(false)
	})

	Required("username", "email", "roles")
})

var CredentialsPayload = Type("Credentials", func() {
	Description("Username and password credentials")
	Attribute("username", String, "Name of user", func() {
		Pattern("^([a-zA-Z0-9@]{4,30})$")
	})
	Attribute("password", String, "Password of user", func() {
		MinLength(6)
		MaxLength(30)
	})
	Required("username", "password")
})

var EmailPayload = Type("EmailPayload", func() {
	Description("Email payload")
	Attribute("email", String, "Email of user", func() {
		Format("email")
	})
	Required("email")
})

// Swagger UI
var _ = Resource("swagger", func() {
	Description("The API swagger specification")

	Files("swagger.json", "swagger/swagger.json")
	Files("swagger-ui/*filepath", "swagger-ui/dist")
})

var JWT = JWTSecurity("jwt", func() {
	Header("Authorization")
	Scope("api:read", "Read API resource")
	Scope("api:write", "Write API resource")
})
