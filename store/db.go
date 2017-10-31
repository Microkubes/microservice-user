package store

import (
	"reflect"
	"time"

	"github.com/JormungandrK/user-microservice/app"
	"github.com/goadesign/goa"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection is an interface to access to the collection struct.
type Collection interface {
	CreateUser(payload *app.UserPayload) (*string, error)
	UpdateUser(userID string, payload *app.UserPayload) (*app.Users, error)
	FindByID(userID string, mediaType *app.Users) error
	FindByUsernameAndPassword(username, password string) (*app.Users, error)
	FindByEmail(email string) (*app.Users, error)
}

// MongoCollection wraps a mgo.Collection to embed methods in models.
type MongoCollection struct {
	*mgo.Collection
}

// NewSession returns a new Mongo Session.
func NewSession(Host string, Username string, Password string, Database string) *mgo.Session {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{Host},
		Username: Username,
		Password: Password,
		Database: Database,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// SetMode - consistency mode for the session.
	session.SetMode(mgo.Monotonic, true)

	return session
}

// PrepareDB ensure presence of persistent and immutable data in the DB.
func PrepareDB(session *mgo.Session, db string, dbCollection string, indexes []string) *mgo.Collection {
	// Create collection
	collection := session.DB(db).C(dbCollection)

	// Define indexes
	for _, elem := range indexes {
		i := []string{elem}
		index := mgo.Index{
			Key:        i,
			Unique:     true,
			DropDups:   true,
			Background: true,
			Sparse:     true,
		}

		// Create indexes
		if err := collection.EnsureIndex(index); err != nil {
			panic(err)
		}
	}

	return collection
}

// CreateUser creates a user if payload is valid, otherwise it returns error
func (c *MongoCollection) CreateUser(payload *app.UserPayload) (*string, error) {
	if payload.Password == nil && payload.ExternalID == nil {
		return nil, goa.ErrBadRequest("password or externalID must be specified!")
	}

	if payload.Password != nil {
		// Hashing password
		userPassword := *payload.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, goa.ErrInternal(err)
		}
		pass := string(hashedPassword)
		payload.Password = &pass
	}

	// Insert Data
	id := bson.NewObjectId()
	err := c.Insert(bson.M{
		"_id":        id,
		"username":   payload.Username,
		"email":      payload.Email,
		"password":   payload.Password,
		"externalId": payload.ExternalID,
		"roles":      payload.Roles,
		"active":     payload.Active,
	})

	// Handle errors
	if err != nil {
		if mgo.IsDup(err) {
			return nil, goa.ErrBadRequest("email or username already exists in the database")
		}
		return nil, goa.ErrInternal(err)
	}

	userID := id.Hex()

	return &userID, nil
}

// UpdateUser updates a user if payload is valid, otherwise it returns error
func (c *MongoCollection) UpdateUser(userID string, payload *app.UserPayload) (*app.Users, error) {
	objectID, err := hexToObjectID(userID)
	if err != nil {
		return nil, err
	}

	err = c.Update(
		bson.M{"_id": objectID},
		bson.M{"$set": payload},
	)

	if err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound(err)
		}
		return nil, goa.ErrInternal(err)
	}

	res := &app.Users{}

	if err = c.FindByID(userID, res); err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID collection by Id in hex representation - real database
func (c *MongoCollection) FindByID(userID string, mediaType *app.Users) error {
	objectID, err := hexToObjectID(userID)
	if err != nil {
		return err
	}

	// Return one user by id.
	if err := c.FindId(objectID).One(&mediaType); err != nil {
		if err.Error() == "not found" {
			return goa.ErrNotFound("user not found")
		}
		return goa.ErrInternal(err)
	}

	return nil
}

// FindByUsernameAndPassword looks up a user by its username and password.
// This is used primarily by other microservices to validate user credentials.
func (c *MongoCollection) FindByUsernameAndPassword(username, password string) (*app.Users, error) {
	query := bson.M{"username": bson.M{"$eq": username}}

	userData := map[string]interface{}{}
	err := c.Collection.Find(query).Limit(1).One(userData)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		print(reflect.TypeOf(err))
		return nil, err
	}
	if _, ok := userData["username"]; !ok {
		return nil, nil
	}
	if _, ok := userData["password"]; !ok {
		return nil, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData["password"].(string)), []byte(password))

	if err != nil {
		return nil, nil
	}
	active, _ := userData["active"].(bool)
	roles := []string{}
	if rolesArr, ok := userData["roles"].([]interface{}); ok {
		for _, role := range rolesArr {
			if roleStr, isString := role.(string); isString {
				roles = append(roles, roleStr)
			}
		}
	}

	if userData["externalId"] == nil {
		userData["externalId"] = ""
	}

	user := &app.Users{
		Active:     active,
		Email:      userData["email"].(string),
		ID:         userData["_id"].(bson.ObjectId).Hex(),
		Roles:      roles,
		ExternalID: userData["externalId"].(string),
		Username:   userData["username"].(string),
	}
	return user, nil
}

// FindByEmail looks up a user by its email.
func (c *MongoCollection) FindByEmail(email string) (*app.Users, error) {
	query := bson.M{"email": bson.M{"$eq": email}}

	userData := map[string]interface{}{}
	err := c.Collection.Find(query).Limit(1).One(userData)
	if err != nil {
		if err.Error() == "not found" {
			return nil, goa.ErrNotFound("user not found")
		}
		return nil, goa.ErrInternal(err)
	}

	active, _ := userData["active"].(bool)
	roles := []string{}
	if rolesArr, ok := userData["roles"].([]interface{}); ok {
		for _, role := range rolesArr {
			if roleStr, isString := role.(string); isString {
				roles = append(roles, roleStr)
			}
		}
	}

	user := &app.Users{
		Active:     active,
		Email:      userData["email"].(string),
		ID:         userData["_id"].(bson.ObjectId).Hex(),
		Roles:      roles,
		ExternalID: userData["externalId"].(string),
		Username:   userData["username"].(string),
	}

	return user, nil
}

// Convert hex representation of object id to bson object id
func hexToObjectID(hexID string) (bson.ObjectId, error) {
	// Return whether userID is a valid hex representation of object id.
	if bson.IsObjectIdHex(hexID) != true {
		return "", goa.ErrBadRequest("invalid user ID")
	}

	// Return an ObjectId from the provided hex representation.
	objectID := bson.ObjectIdHex(hexID)

	// Return true if objectID is valid. A valid objectID must contain exactly 12 bytes.
	if objectID.Valid() != true {
		return "", goa.ErrInternal("invalid object ID")
	}

	return objectID, nil
}
