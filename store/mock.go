package store

import (
	"sync"
	"errors"
	// "reflect"

	"user-microservice/app"
	"gopkg.in/mgo.v2/bson"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	users		map[string]*app.UserPayload
}

// NewDB initializes a new "DB" with dummy data.
func NewDB() *DB {
	roles := []string{"admin", "user"}
	user := &app.UserPayload{
		Active: false,
		Email: "frieda@oberbrunnerkirlin.name",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		Roles: roles,
		Username: "User1",
		Password: "pass",
	}
	return &DB{users: map[string]*app.UserPayload{"5975c461f9f8eb02aae053f3": user}}
}

// Reset removes all entries from the database.
func (db *DB) Reset() {
	db.users = make(map[string]*app.UserPayload)
}

// Mock implementation
func (db *DB) FindByID(objectId bson.ObjectId, mediaType *app.Users) error {
	db.Lock()
	defer db.Unlock()

	id := objectId.Hex()

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

// Insert mock.
func (db *DB) Insert(docs ...interface{}) error {
	return nil
}
