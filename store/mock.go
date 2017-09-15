package store

import (
	"fmt"
	"sync"

	"github.com/JormungandrK/user-microservice/app"
	"github.com/goadesign/goa"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	users map[string]*app.UserPayload
}

// NewDB initializes a new "DB" with dummy data.
func NewDB() *DB {
	roles := []string{"admin", "user"}
	pass := "pass"
	extID := "qwerc461f9f8eb02aae053f3"
	user := &app.UserPayload{
		Active:     false,
		Email:      "frieda@oberbrunnerkirlin.name",
		ExternalID: &extID,
		Roles:      roles,
		Username:   "User1",
		Password:   &pass,
	}
	return &DB{users: map[string]*app.UserPayload{"5975c461f9f8eb02aae053f3": user}}
}

// Reset removes all entries from the database.
func (db *DB) Reset() {
	db.users = make(map[string]*app.UserPayload)
}

// FindByID mock implementation
func (db *DB) FindByID(userID string, mediaType *app.Users) error {
	db.Lock()
	defer db.Unlock()

	if userID == "6975c461f9f8eb02aae053f4" {
		return goa.ErrInternal("internal server error")
	}

	if userID == "fakeobjectidab02aae053f3aasadas" {
		return goa.ErrBadRequest("invalid user ID")
	}

	if user, ok := db.users[userID]; ok {
		mediaType.ID = userID
		mediaType.Username = user.Username
		mediaType.Email = user.Email
		mediaType.ExternalID = *user.ExternalID
		mediaType.Roles = user.Roles
		mediaType.Active = user.Active
	} else {
		return goa.ErrNotFound("user not found")
	}

	return nil
}

// Insert mock implementation
func (db *DB) CreateUser(payload *app.UserPayload) (*string, error) {
	if payload.Password == nil && payload.ExternalID == nil {
		return nil, goa.ErrBadRequest("password or externalID must be specified!")
	}
	if payload.Username == "internal-error-user" {
		return nil, goa.ErrInternal("internal server error")
	}

	id := "ab75c461f9f8eb02aae058zr"
	db.users[id] = payload

	return &id, nil
}

// Update mock implementation
func (db *DB) UpdateUser(userID string, payload *app.UserPayload) (*app.Users, error) {
	if userID == "6975c461f9f8eb02aae053f4" {
		return nil, goa.ErrInternal("internal server error")
	}
	if userID == "fakeobjectidab02aae053f3aasadas" {
		return nil, goa.ErrBadRequest("invalid user ID")
	}

	if _, ok := db.users[userID]; ok {
		db.users[userID] = payload
		user, _ := db.users[userID]
		return &app.Users{
			Active:     user.Active,
			Email:      user.Email,
			ExternalID: *user.ExternalID,
			ID:         userID,
			Roles:      user.Roles,
			Username:   user.Username,
		}, nil
	} else {
		return nil, goa.ErrNotFound("user not found")
	}

	return nil, nil
}

func (db *DB) FindByUsernameAndPassword(username, password string) (*app.Users, error) {
	if username == "validuser" && password == "valid-pass" {
		return &app.Users{
			Active:     true,
			Email:      "email@example.com",
			ExternalID: "1234",
			ID:         "000000000000001",
			Roles:      []string{"user"},
			Username:   "validuser",
		}, nil
	}
	if username == "internal-error-user" {
		return nil, fmt.Errorf("Internal server error")
	}
	return nil, nil
}

func (db *DB) FindByEmail(email string) (*app.Users, error) {
	if email == "frieda@oberbrunnerkirlin.name" {
		return &app.Users{
			Active:     true,
			Email:      "frieda@oberbrunnerkirlin.name",
			ExternalID: "qwerc461f9f8eb02aae053f3",
			ID:         "5975c461f9f8eb02aae053f3",
			Roles:      []string{"admin", "user"},
			Username:   "User1",
		}, nil
	}
	if email == "example@notexists.com" {
		return nil, goa.ErrNotFound("user not exist")
	}
	if email == "example@invalid.com" {
		return nil, goa.ErrInternal("internal server error")
	}

	return nil, nil
}
