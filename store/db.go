package store

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"user-microservice/app"
)

const (
    Host     = "127.0.0.1:27017"
    Username = "restapi"
    Password = "restapi"
    Database = "users"
)

type Collection interface {
	FindByID(id string, mediaType *app.Users) error
}

type MongoCollection struct {
	*mgo.Collection
}

// NewSession returns a new Mongo Session.
func NewSession() *mgo.Session {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{Host},
		Username: Username,
		Password: Password,
		Database: Database,
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
	index := mgo.Index{
		Key:        indexes,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	// Create indexes
	if err := collection.EnsureIndex(index); err != nil {
		panic(err)
	}

	return collection
}

// Find collection by Id in hex representation - real database
func (c *MongoCollection) FindByID(id string, mediaType *app.Users) error {
	// Return whether ctx.UserID is a valid hex representation of an ObjectId.
	if bson.IsObjectIdHex(id) != true {
		err := errors.New("Invalid Id")
		return err
	}

	// Return an ObjectId from the provided hex representation. 
    userId := bson.ObjectIdHex(id)

	// Return true if userId is valid. A valid userId must contain exactly 12 bytes.
	if userId.Valid() != true {
		err := errors.New("Invalid Id")
		return err
	}

	// Return one user by id.
	if err := c.FindId(userId).One(&mediaType); err != nil {
		return err
	}

	return nil
}
