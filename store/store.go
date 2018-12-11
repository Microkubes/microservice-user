package store

import (
	"github.com/Microkubes/backends"
	"github.com/Microkubes/microservice-user/app"
)

type UserRecord struct {
	ID string `json:"id" bson:"_id"`
	// Status of user account
	Active bool `form:"active" json:"active" yaml:"active" xml:"active"`
	// Email of user
	Email string `form:"email" json:"email" yaml:"email" xml:"email"`
	// External id of user
	ExternalID string `form:"externalId,omitempty" json:"externalId,omitempty" yaml:"externalId,omitempty" xml:"externalId,omitempty"`
	// List of namespaces this user belongs to
	Namespaces []string `form:"namespaces,omitempty" json:"namespaces,omitempty" yaml:"namespaces,omitempty" xml:"namespaces,omitempty"`
	// List of organizations to which this user belongs to
	Organizations []string `form:"organizations,omitempty" json:"organizations,omitempty" yaml:"organizations,omitempty" xml:"organizations,omitempty"`
	// Password of user
	Password string `form:"password,omitempty" json:"password,omitempty" yaml:"password,omitempty" xml:"password,omitempty"`
	// Roles of user
	Roles []string `form:"roles,omitempty" json:"roles,omitempty" yaml:"roles,omitempty" xml:"roles,omitempty"`
	// Token for email verification
	Token string `form:"token,omitempty" json:"token,omitempty" yaml:"token,omitempty" xml:"token,omitempty"`
}

func (u *UserRecord) ToAppUsers() *app.Users {
	au := &app.Users{
		Active:        u.Active,
		Email:         u.Email,
		ExternalID:    u.ExternalID,
		ID:            u.ID,
		Namespaces:    u.Namespaces,
		Organizations: u.Organizations,
		Roles:         u.Roles,
	}
	return au
}

// User wraps User's collections/tables. Implements backneds.Repository interface
type User struct {
	Users  backends.Repository
	Tokens backends.Repository
}
