package store

import (
	"errors"
	"fmt"
	"sync"

	"github.com/JormungandrK/user-microservice/app"

	"gopkg.in/mgo.v2/bson"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	users map[string]*app.UserPayload
}

// NewDB initializes a new "DB" with dummy data.
func NewDB() *DB {
	roles := []string{"admin", "user"}
	user := &app.UserPayload{
		Active:     false,
		Email:      "frieda@oberbrunnerkirlin.name",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		Roles:      roles,
		Username:   "User1",
		Password:   "pass",
	}
	return &DB{users: map[string]*app.UserPayload{"5975c461f9f8eb02aae053f3": user}}
}

// Reset removes all entries from the database.
func (db *DB) Reset() {
	db.users = make(map[string]*app.UserPayload)
}

// FindByID mock implementation
func (db *DB) FindByID(objectID bson.ObjectId, mediaType *app.Users) error {
	db.Lock()
	defer db.Unlock()

	id := objectID.Hex()

	if user, ok := db.users[id]; ok {
		mediaType.ID = id
		mediaType.Username = user.Username
		mediaType.Email = user.Email
		mediaType.ExternalID = user.ExternalID
		mediaType.Roles = user.Roles
		mediaType.Active = user.Active
	} else {
		err := errors.New("Cannot retrieve collection by Id")
		return err
	}

	return nil
}

// Insert mock implementation
func (db *DB) Insert(docs ...interface{}) error {
	return nil
}

// Update mock implementation
func (db *DB) Update(selector interface{}, update interface{}) error {
	return nil
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
