package store

import (
	"sync"
	// "errors"
	"reflect"

	"user-microservice/app"
	"gopkg.in/mgo.v2/bson"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	users		map[string]interface{}
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
	return &DB{users: map[string]interface{}{"5975c461f9f8eb02aae053f3": user}}
}

// Reset removes all entries from the database.
func (db *DB) Reset() {
	db.users = make(map[string]interface{})
}

// Mock implementation
func (db *DB) FindByID(objectId bson.ObjectId, mediaType *app.Users) error {
	db.Lock()
	defer db.Unlock()

	id := objectId.Hex()

	if user, ok := db.users[id]; ok {
		s := reflect.ValueOf(&user)
		typeOfUser := s.Type()
		// for i := 0; i < s.NumField(); i++ {
		// 	f := s.Field(i)

		// 	fmt.Printf("%d: %s %s = %v\n", i,
		// 		typeOfT.Field(i).Name, f.Type(), f.Interface())
		// }
		mediaType.ID = id
		mediaType.Active = s.FieldByName("Active").(bool)
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

	db.users["3375c461f9f8eb02aae053q4"] = docs				

	return nil
}
