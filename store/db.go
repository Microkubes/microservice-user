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

// IUserCollection is an interface to access the userCollection struct.
type IUserCollection interface {
	CreateUser(payload *app.UserPayload) (*string, error)
	UpdateUser(userID string, payload *app.UserPayload) (*app.Users, error)
	FindByID(userID string, mediaType *app.Users) error
	FindByEmailAndPassword(email, password string) (*app.Users, error)
	FindByEmail(email string) (*app.Users, error)
	ActivateUser(email string) error
}

// ITokenCollection is an interface to access the tokenCollection struct.
type ITokenCollection interface {
	CreateToken(payload *app.UserPayload) error
	VerifyToken(token string) (*string, error)
	DeleteToken(token string) error
}

// UserCollection wraps a mgo.Collection to embed methods in models.
type UserCollection struct {
	*mgo.Collection
}

// TokenCollection wraps a mgo.Collection to embed methods in models.
type TokenCollection struct {
	*mgo.Collection
}

type Collections struct {
	Users  IUserCollection
	Tokens ITokenCollection
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
func PrepareDB(session *mgo.Session, db string, dbCollection string, indexes []string, enableTTL bool) *mgo.Collection {
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

	if enableTTL == true {
		index := mgo.Index{
			Key:         []string{"created_at"},
			Unique:      false,
			DropDups:    false,
			Background:  true,
			Sparse:      true,
			ExpireAfter: time.Duration(60*60*24) * time.Second,
		}
		if err := collection.EnsureIndex(index); err != nil {
			panic(err)
		}

	}

	return collection
}

// CreateUser creates a user if payload is valid, otherwise it returns error
func (c *UserCollection) CreateUser(payload *app.UserPayload) (*string, error) {
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
		"email":      payload.Email,
		"password":   payload.Password,
		"externalId": payload.ExternalID,
		"roles":      payload.Roles,
		"active":     payload.Active,
	})

	// Handle errors
	if err != nil {
		if mgo.IsDup(err) {
			return nil, goa.ErrBadRequest("email already exists in the database")
		}
		return nil, goa.ErrInternal(err)
	}

	userID := id.Hex()

	return &userID, nil
}

// UpdateUser updates a user if payload is valid, otherwise it returns error
func (c *UserCollection) UpdateUser(userID string, payload *app.UserPayload) (*app.Users, error) {
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
func (c *UserCollection) FindByID(userID string, mediaType *app.Users) error {
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

// FindByEmailAndPassword looks up a user by its email and password.
// This is used primarily by other microservices to validate user credentials.
func (c *UserCollection) FindByEmailAndPassword(email, password string) (*app.Users, error) {
	query := bson.M{"email": bson.M{"$eq": email}}

	userData := map[string]interface{}{}
	err := c.Collection.Find(query).Limit(1).One(userData)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		print(reflect.TypeOf(err))
		return nil, err
	}
	if _, ok := userData["email"]; !ok {
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

	var externalID string
	if userData["externalId"] == nil {
		externalID = ""
	} else {
		externalID = userData["externalId"].(string)
	}

	user := &app.Users{
		Active:     active,
		Email:      userData["email"].(string),
		ID:         userData["_id"].(bson.ObjectId).Hex(),
		Roles:      roles,
		ExternalID: externalID,
	}
	return user, nil
}

// FindByEmail looks up a user by its email.
func (c *UserCollection) FindByEmail(email string) (*app.Users, error) {
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
	}

	return user, nil
}

func (c *UserCollection) ActivateUser(email string) error {
	err := c.Update(
		bson.M{"email": email},
		bson.M{"$set": bson.M{"active": true}},
	)
	if err != nil {
		if err.Error() == "not found" {
			return goa.ErrNotFound(err)
		}
		return goa.ErrInternal(err)
	}
	return nil
}

func (c *TokenCollection) CreateToken(payload *app.UserPayload) error {
	id := bson.NewObjectId()
	err := c.Insert(bson.M{
		"_id":        id,
		"email":      payload.Email,
		"token":      payload.Token,
		"created_at": time.Now(),
	})

	if err != nil {
		return goa.ErrInternal(err)
	}
	return nil
}

func (c *TokenCollection) VerifyToken(token string) (*string, error) {
	query := bson.M{"token": bson.M{"$eq": token}}

	tokenData := map[string]interface{}{}
	err := c.Collection.Find(query).Limit(1).One(tokenData)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, goa.ErrNotFound("token not found!")
		}
		print(reflect.TypeOf(err))
		return nil, goa.ErrInternal(err)
	}

	email := tokenData["email"].(string)
	return &email, nil
}

func (c *TokenCollection) DeleteToken(token string) error {
	query := bson.M{"token": bson.M{"$eq": token}}

	err := c.Collection.Remove(query)
	if err != nil {
		if err == mgo.ErrNotFound {
			return goa.ErrNotFound("token not found!")
		}
		print(reflect.TypeOf(err))
		return goa.ErrInternal(err)
	}
	return nil
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
