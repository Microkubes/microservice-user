package store

import (
	"sync"

	"github.com/Microkubes/backends"
	"gopkg.in/mgo.v2/bson"
)

type MapStore map[string]interface{}

// DB emulatesl a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	MapStore
}

var (
	BAD_REQUEST    = "bad request"
	NOT_FOUND      = "not found"
	INTERNAL_ERROR = "internal server error"
)

// NewDB initializes a new "DB" with dummy data.
func NewDB() User {

	users := &DB{
		MapStore: map[string]interface{}{
			"5df2103b5f1b640001142d3c": map[string]interface{}{
				"id":         "5df2103b5f1b640001142d3c",
				"email":      "keitaro-user1@gmail.com",
				"password":   "keitaro",
				"externalID": "some-id",
				"roles":      []string{"user"},
			},
		},
	}

	tokens := &DB{
		MapStore: map[string]interface{}{
			"z8cfa84f-bb6c-4c84-b39b-76dd32653999": map[string]interface{}{
				"id":    "z8cfa84f-bb6c-4c84-b39b-76dd32653999",
				"email": "keitaro-user1@gmail.com",
				"token": "sdaewefdc234erfdd123erfdxc23edx",
			},
		},
	}

	return User{
		Users:  users,
		Tokens: tokens,
	}
}

func (db *DB) GetOne(filter backends.Filter, result interface{}) (interface{}, error) {

	db.Lock()
	defer db.Unlock()

	if id, ok := filter["id"]; ok {
		idString := id.(string)

		if idString == "bad-id" {
			return nil, backends.ErrInvalidInput(BAD_REQUEST)
		}

		if idString == "internal-err-id" {
			return nil, backends.ErrBackendError(INTERNAL_ERROR)
		}

		record, ok := db.MapStore[idString]
		if !ok {
			return nil, backends.ErrNotFound(NOT_FOUND)
		}

		err := backends.MapToInterface(&record, &result)
		if err != nil {
			return nil, backends.ErrBackendError(err)
		}
	}

	if email, ok := filter["email"]; ok {
		emailString := email.(string)

		if emailString == "internal-error@example.com" {
			return nil, backends.ErrBackendError(INTERNAL_ERROR)
		}

		if emailString == "not-found@gmail.com" {
			return nil, backends.ErrNotFound(NOT_FOUND)
		}

		if emailString == "bad@gmail.com" {
			return nil, backends.ErrInvalidInput(BAD_REQUEST)
		}

		for _, r := range db.MapStore {
			record := r.(map[string]interface{})

			if record["email"] == emailString {
				err := backends.MapToInterface(record, &result)
				if err != nil {
					return nil, backends.ErrBackendError(err)
				}

				break
			}
		}
	}

	if token, ok := filter["token"]; ok {
		tokenString := token.(string)

		if tokenString == "internal-error-token" {
			return nil, backends.ErrBackendError(INTERNAL_ERROR)
		}

		if tokenString == "not-found-token" {
			return nil, backends.ErrNotFound(NOT_FOUND)
		}

		for _, r := range db.MapStore {
			record := r.(map[string]interface{})

			if record["token"] == tokenString {
				err := backends.MapToInterface(record, &result)
				if err != nil {
					return nil, backends.ErrBackendError(err)
				}

				break
			}
		}
	}

	return result, nil
}

func (db *DB) GetAll(filter backends.Filter, results interface{}, order string, sorting string, limit int, offset int) (interface{}, error) {
	if offset > 2 {
		return nil, backends.ErrNotFound("offset is too high")
	}

	var users []map[string]interface{}
	for _, v := range db.MapStore {
		users = append(users, v.(map[string]interface{}))
	}

	if len(users) == 0 {
		return nil, backends.ErrNotFound("Empty users")
	}

	return users, nil
}

func (db *DB) Save(object interface{}, filter backends.Filter) (interface{}, error) {

	db.Lock()
	defer db.Unlock()

	var result interface{}

	payload, err := backends.InterfaceToMap(object)
	if err != nil {
		return nil, backends.ErrBackendError(err)
	}

	if filter == nil {

		if (*payload)["email"] == "internal-error@example.com" {
			return nil, backends.ErrBackendError(INTERNAL_ERROR)
		}

		id := bson.NewObjectId().Hex()

		(*payload)["id"] = id

		db.MapStore[id] = *payload
		result = &UserRecord{
			ExternalID: "",
			Roles:      []string{},
		}
	} else {

		if id, ok := filter["id"]; ok {
			idString := id.(string)

			if idString == "bad-id" {
				return nil, backends.ErrInvalidInput(BAD_REQUEST)
			}

			if idString == "internal-err-id" {
				return nil, backends.ErrBackendError(INTERNAL_ERROR)
			}

			record, ok := db.MapStore[idString]
			if !ok {
				return nil, backends.ErrNotFound(NOT_FOUND)
			}

			updateRecord := record.(map[string]interface{})
			for k, v := range *payload {
				updateRecord[k] = v
			}

			payload = &updateRecord
		}

		if email, ok := filter["email"]; ok {
			emailString := email.(string)

			if emailString == "bad-id" {
				return nil, backends.ErrInvalidInput(BAD_REQUEST)
			}

			if emailString == "internal-err-id" {
				return nil, backends.ErrBackendError(INTERNAL_ERROR)
			}

			for _, r := range db.MapStore {
				record := r.(map[string]interface{})

				if record["email"] == emailString {

					payload = &record
					break
				}
			}
		}
	}

	err = backends.MapToInterface(payload, &result)
	if err != nil {
		return nil, backends.ErrBackendError(err)
	}

	return result, nil
}

func (db *DB) DeleteOne(filter backends.Filter) error {

	db.Lock()
	defer db.Unlock()

	if token, ok := filter["token"]; ok {
		tokenString := token.(string)

		for key, r := range db.MapStore {
			record := r.(map[string]interface{})

			if record["token"] == tokenString {

				delete(db.MapStore, key)
				break
			}
		}
	}

	return nil
}

func (db *DB) DeleteAll(filter backends.Filter) error {
	return nil
}
