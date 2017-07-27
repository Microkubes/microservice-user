package store

import (
	"sync"
	"errors"
	"user-microservice/app"
	"gopkg.in/mgo.v2/bson"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	users		map[string]*UserModel
}

// UserModel is the database "users" for users
type UserModel struct {
	Active bool 
	Email string 
	ExternalID string 
	ID string
	Roles []string 
	Username string
	Password string
}

// NewDB initializes a new "DB" with dummy data.
func NewDB() *DB {
	roles := []string{"admin", "user"}
	user := &UserModel{
		Active: false,
		Email: "frieda@oberbrunnerkirlin.name",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		ID: "5975c461f9f8eb02aae053f3",
		Roles: roles,
		Username: "User1",
		Password: "pass",
	}
	return &DB{users: map[string]*UserModel{"5975c461f9f8eb02aae053f3": user}}
}

// Reset removes all entries from the database.
func (db *DB) Reset() {
	db.users = make(map[string]*UserModel)
}

// Mock implementation
func (db *DB) FindByID(objectId bson.ObjectId, mediaType *app.Users) error {
	db.Lock()
	defer db.Unlock()

	id := objectId.Hex()

	if user, ok := db.users[id]; ok {
		mediaType.ID = user.ID
		mediaType.Active = user.Active
		mediaType.Email = user.Email
		mediaType.ExternalID = user.ExternalID
		mediaType.Roles = user.Roles
		mediaType.Username = user.Username
	} else {
		err := errors.New("Cannot retrieve collection by Id")
		return err				
	}

	return nil
}

func (db *DB) Insert(docs ...interface{}) error {
	db.Lock()
	defer db.Unlock()


	// email := docs.(struct{Email string}).Email
	roles := []string{"admin", "user"}

	user := &UserModel{
        Active: false,
        Email: "email@gmail.com",
        ExternalID: "exidnew5975c461f9f8eb02aae053f3",
        ID: "new5975c461f9f8eb02aae053f3",
        Roles: roles,
        Username: "user",
		Password: "pass",
    }
    
	db.users["5975c461f9f8eb02aae053f4"] = user
	
	return nil
}
