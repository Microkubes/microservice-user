package store

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"user-microservice/app"
)

type Collection interface {
	FindByID(id bson.ObjectId, mediaType *app.Users) error
	Insert(docs ...interface{}) error
}

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
func (c *MongoCollection) FindByID(objectId bson.ObjectId, mediaType *app.Users) error {
	// Return one user by id.
	if err := c.FindId(objectId).One(&mediaType); err != nil {
		return err
	}

	return nil
}
