package store

import (
	"reflect"
	"time"

	"github.com/JormungandrK/user-microservice/app"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Collection is an interface to access to the collection struct.
type Collection interface {
	FindByID(id bson.ObjectId, mediaType *app.Users) error
	Insert(docs ...interface{}) error
	Update(selector interface{}, update interface{}) error
	FindByUsernameAndPassword(username, password string) (*app.Users, error)
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

// FindByID collection by Id in hex representation - real database
func (c *MongoCollection) FindByID(objectID bson.ObjectId, mediaType *app.Users) error {
	// Return one user by id.
	if err := c.FindId(objectID).One(&mediaType); err != nil {
		return err
	}

	return nil
}

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
