package design

// Use . imports to enable the DSL
import (
	. "github.com/keitaroinc/goa/design"
	. "github.com/keitaroinc/goa/design/apidsl"
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
	BasePath("/")
	// Do not setup security here!

	// Allow preflight requests (HTTP OPTIONS)
	Origin("*", func() {
		Methods("OPTIONS")
	})

	// Actions define a single API endpoint
	Action("create", func() {
		Description("Creates user")
		Routing(POST(""))
		Payload(CreateUserPayload)
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
		Response(OK, UserMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getMe", func() {
		Description("Retrieves the user information for the authenticated user")
		Routing(GET("/me"))
		Response(OK, UserMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("getAll", func() {
		Description("Retrieves all active users")
		Routing(GET(""))
		Params(func() {
			Param("order", String, "Order by")
			Param("sorting", String, func() {
				Enum("asc", "desc")
			})
			Param("limit", Integer, "Limit users per page")
			Param("offset", Integer, "Number of users to skip")
		})
		Response(OK)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("update", func() {
		Description("Update user")
		Routing(PUT("/:userId"))
		Params(func() {
			Param("userId", String, "User ID")
		})
		Payload(UpdateUserPayload)
		Response(OK, UserMedia)
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("findUsers", func() {
		Description("Find (filter) users by some filter.")
		Routing(POST("/list"))
		Payload(FilterPayload)
		Response(OK, UsersPageMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})

	Action("find", func() {
		Description("Find a user by email+password")
		Routing(POST("find"))
		Payload(CredentialsPayload)
		Response(OK, UserMedia)
		Response(NotFound, ErrorMedia)
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

	Action("verify", func() {
		Description("Verify a user by token")
		Routing(GET("verify"))
		Params(func() {
			Param("token", String, "Token")
		})
		Response("OK", func() {
			Description("User is verified")
			Status(200)
		})
		Response(NotFound, ErrorMedia)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
	Action("resetVerificationToken", func() {
		Description("Reset verification token")
		Routing(POST("verification/reset"))
		Payload(EmailPayload)
		Response("OK", func() {
			Description("Verification token reset")
			Status(200)
			Media(ResetTokenMedia)
		})
		Response(BadRequest, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
	Action("forgotPassword", func() {
		Description("Forgot password action (sending email to user with link for resseting password)")
		Routing(POST("password/forgot"))
		Payload(EmailPayload)
		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
	Action("forgotPasswordUpdate", func() {
		Description("Password token validation & password update")
		Routing(PUT("password/forgot"))
		Payload(ForgotPasswordPayload)
		Response(OK)
		Response(BadRequest, ErrorMedia)
		Response(NotFound, ErrorMedia)
		Response(InternalServerError, ErrorMedia)
	})
})

// UserMedia defines the media type used to render user.
var UserMedia = MediaType("application/vnd.goa.user+json", func() {
	TypeName("users")
	Reference(CreateUserPayload)

	Attributes(func() {
		Attribute("id", String, "Unique user ID")
		Attribute("email")
		Attribute("roles")
		Attribute("externalId")
		Attribute("active")
		Attribute("organizations")
		Attribute("namespaces")
		Required("id", "email", "roles", "externalId", "active")
	})

	View("default", func() {
		Attribute("id")
		Attribute("email")
		Attribute("roles")
		Attribute("externalId")
		Attribute("active")
		Attribute("organizations")
		Attribute("namespaces")
	})
})

// ResetTokenMedia is returned after successful reset of the verification token
var ResetTokenMedia = MediaType("ResetTokenMedia", func() {
	TypeName("ResetToken")
	Attributes(func() {
		Attribute("id", String, "User ID")
		Attribute("email", String, "User email")
		Attribute("token", String, "New token")
		Required("id", "email", "token")
	})
	View("default", func() {
		Attribute("id")
		Attribute("email")
		Attribute("token")
	})
})

// CreateUserPayload defines the payload for the user.
var CreateUserPayload = Type("CreateUserPayload", func() {
	Description("CreateUserPayload")

	Attribute("email", String, "Email of user", func() {
		Format("email")
	})
	Attribute("password", String, "Password of user", func() {
		MinLength(6)
		MaxLength(30)
	})
	Attribute("roles", ArrayOf(String), "Roles of user")
	Attribute("organizations", ArrayOf(String), "List of organizations to which this user belongs to")
	Attribute("namespaces", ArrayOf(String), "List of namespaces this user belongs to")
	Attribute("externalId", String, "External id of user")
	Attribute("active", Boolean, "Status of user account", func() {
		Default(false)
	})
	Attribute("token", String, "Token for email verification")

	Required("email")
})

// UpdateUserPayload defines the payload for the user.
var UpdateUserPayload = Type("UpdateUserPayload", func() {
	Description("UpdateUserPayload")

	Attribute("email", String, "Email of user", func() {
		Format("email")
	})
	Attribute("password", String, "Password of user", func() {
		MinLength(6)
		MaxLength(30)
	})
	Attribute("roles", ArrayOf(String), "Roles of user")
	Attribute("organizations", ArrayOf(String), "List of organizations to which this user belongs to")
	Attribute("namespaces", ArrayOf(String), "List of namespaces this user belongs to")
	Attribute("externalId", String, "External id of user")
	Attribute("active", Boolean, "Status of user account", func() {
		Default(false)
	})
	Attribute("token", String, "Token for email verification")
})

// CredentialsPayload defines the payload for the credentials.
var CredentialsPayload = Type("Credentials", func() {
	Description("Email and password credentials")
	Attribute("email", String, "Email of user", func() {
		Format("email")
	})
	Attribute("password", String, "Password of user", func() {
		MinLength(6)
		MaxLength(30)
	})
	Required("email", "password")
})

// EmailPayload defines the payload for the email.
var EmailPayload = Type("EmailPayload", func() {
	Description("Email payload")
	Attribute("email", String, "Email of user", func() {
		Format("email")
	})
	Required("email")
})

// ForgotPasswordPayload defines the payload for the password/forgot.
var ForgotPasswordPayload = Type("ForgotPasswordPayload", func() {
	Description("Password Reset payload")
	Attribute("email", String, "Email of the user", func() {
		Format("email")
	})
	Attribute("password", String, "New password", func() {
		MinLength(6)
		MaxLength(30)
	})
	Attribute("token", String, "Forgot password token")
	Required("password", "token")
})

// Swagger UI
var _ = Resource("swagger", func() {
	Description("The API swagger specification")

	BasePath("/users")
	Files("swagger.json", "swagger/swagger.json")
	Files("swagger-ui/*filepath", "swagger-ui/dist")
})

// FilterPayload Users filter request payload.
var FilterPayload = Type("FilterPayload", func() {
	Attribute("page", Integer, "Page number (1-based).")
	Attribute("pageSize", Integer, "Items per page.")
	Attribute("filter", ArrayOf(FilterProperty), "Users filter.")
	Attribute("sort", OrderSpec, "Sort specification.")
	Required("page", "pageSize")
})

// FilterProperty Single property filter. Holds the property name and the value to be matched for that property.
var FilterProperty = Type("FilterProperty", func() {
	Attribute("property", String, "Property name")
	Attribute("value", String, "Property value to match")
	Required("property", "value")
})

// OrderSpec specifies the sorting - by which property and the direction, either 'asc' (ascending)
// or 'desc' (descending).
var OrderSpec = Type("OrderSpec", func() {
	Attribute("property", String, "Sort by property")
	Attribute("direction", String, "Sort order. Can be 'asc' or 'desc'.")
	Required("property", "direction")
})

// UsersPageMedia result of filter-by. One result page along with items (array of Users).
var UsersPageMedia = MediaType("application/mt.ckan.users-page+json", func() {
	TypeName("UsersPage")
	Attributes(func() {
		Attribute("page", Integer, "Page number (1-based).")
		Attribute("pageSize", Integer, "Items per page.")
		Attribute("items", ArrayOf(UserMedia), "Users list")
	})
	View("default", func() {
		Attribute("page")
		Attribute("pageSize")
		Attribute("items")
	})
})
