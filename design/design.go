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

        // Actions define a single API endpoint
        Action("create", func() {
            Description("Creates user")
            Routing(POST("/"))
            Payload(UserPayload)
            Response(Created, UserMedia)
        })

        Action("get", func() {
                Description("Get user by id")
                Routing(GET("/:userId"))
                Params(func() {
                        Param("userId", Integer, "User ID")
                })
                Response(OK)
                Response(NotFound)
        })

        Action("getMe", func() {
                Description("Retrieves the user information for the authenticated user")
                Routing(GET("/me"))
                Response(OK)
                Response(NotFound)
        })

        Action("update", func() {
            Description("Update user")
            Routing(PUT("/:userId"))
            Params(func() {
                    Param("userId", Integer, "User ID")
            })
            Payload(UserPayload)
            Response(NotFound)
            Response(OK, UserMedia)
        })
})

// UserMedia defines the media type used to render user.
var UserMedia = MediaType("application/vnd.goa.user+json", func() {
        TypeName("users")
        Reference(UserPayload)

        Attributes(func() {                         
                Attribute("id", Integer, "Unique user ID")
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

        Required("username", "email", "password", "roles", "externalId")
})

// Swagger UI
var _ = Resource("swagger", func() {
        Description("The API swagger specification")

        Files("swagger.json", "swagger/swagger.json")
        Files("swagger-ui/*filepath", "swagger-ui/dist")
})
