package store

import (
	"github.com/JormungandrK/backends"
)

// User wraps User's collections/tables. Implements backneds.Repository interface
type User struct {
	Users  backends.Repository
	Tokens backends.Repository
}
