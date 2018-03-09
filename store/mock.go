package store

import (
	"fmt"
	"sync"

	"github.com/JormungandrK/microservice-user/app"
	"github.com/goadesign/goa"
)

// DB emulates a database driver using in-memory data structures.
type DB struct {
	sync.Mutex
	users map[string]*app.UserPayload
}

// NewDB initializes a new "DB" with dummy data.
func NewDB() Collections {
	roles := []string{"admin", "user"}
	pass := "pass"
	extID := "qwerc461f9f8eb02aae053f3"
	user := &app.UserPayload{
		Active:     false,
		Email:      "frieda@oberbrunnerkirlin.name",
		ExternalID: &extID,
		Roles:      roles,
		Password:   &pass,
	}
	//return &DB{users: map[string]*app.UserPayload{"5975c461f9f8eb02aae053f3": user}}
	tokens := TokensMock{
		Tokens: map[string]*app.UserPayload{},
	}
	users := DB{
		users: map[string]*app.UserPayload{
			"5975c461f9f8eb02aae053f3": user,
			"5975c461f9f8eb02aae053f4": &app.UserPayload{
				Active: false,
				Email:  "email@example.com",
				Roles:  []string{"user"},
			},
			"5975c461f9f8eb02aae053f5": &app.UserPayload{
				Active: true,
				Email:  "already-active@example.com",
				Roles:  []string{"user"},
			},
		},
	}
	colls := Collections{
		Tokens: &tokens,
		Users:  &users,
	}
	return colls
}

// Reset removes all entries from the database.
func (db *DB) Reset() {
	db.users = make(map[string]*app.UserPayload)
}

// FindByID mock implementation
func (db *DB) FindByID(userID string, mediaType *app.Users) error {
	db.Lock()
	defer db.Unlock()

	if userID == "6975c461f9f8eb02aae053f4" {
		return goa.ErrInternal("internal server error")
	}

	if userID == "fakeobjectidab02aae053f3aasadas" {
		return goa.ErrBadRequest("invalid user ID")
	}

	if user, ok := db.users[userID]; ok {
		mediaType.ID = userID
		mediaType.Email = user.Email
		mediaType.ExternalID = *user.ExternalID
		mediaType.Roles = user.Roles
		mediaType.Active = user.Active
	} else {
		return goa.ErrNotFound("user not found")
	}

	return nil
}

// CreateUser mock implementation
func (db *DB) CreateUser(payload *app.UserPayload) (*string, error) {
	if payload.Password == nil && payload.ExternalID == nil {
		return nil, goa.ErrBadRequest("password or externalID must be specified!")
	}
	if payload.Email == "internal-error@example.com" {
		return nil, goa.ErrInternal("internal server error")
	}

	id := "ab75c461f9f8eb02aae058zr"
	db.users[id] = payload

	return &id, nil
}

// UpdateUser mock implementation
func (db *DB) UpdateUser(userID string, payload *app.UserPayload) (*app.Users, error) {
	if userID == "6975c461f9f8eb02aae053f4" {
		return nil, goa.ErrInternal("internal server error")
	}
	if userID == "fakeobjectidab02aae053f3aasadas" {
		return nil, goa.ErrBadRequest("invalid user ID")
	}

	if _, ok := db.users[userID]; ok {
		db.users[userID] = payload
		user, _ := db.users[userID]
		return &app.Users{
			Active:     user.Active,
			Email:      user.Email,
			ExternalID: *user.ExternalID,
			ID:         userID,
			Roles:      user.Roles,
		}, nil
	}
	return nil, goa.ErrNotFound("user not found")
}

// FindByEmailAndPassword mock implementation
func (db *DB) FindByEmailAndPassword(email, password string) (*app.Users, error) {
	if email == "email@example.com" && password == "valid-pass" {
		externalID := "1234"
		return &app.Users{
			Active:     true,
			Email:      "email@example.com",
			ExternalID: externalID,
			ID:         "000000000000001",
			Roles:      []string{"user"},
		}, nil
	}
	if email == "internal-error@example.com" {
		return nil, fmt.Errorf("Internal server error")
	}
	return nil, nil
}

// FindByEmail mock implementation
func (db *DB) FindByEmail(email string) (*app.Users, error) {
	if email == "frieda@oberbrunnerkirlin.name" {
		externalID := "qwerc461f9f8eb02aae053f3"
		return &app.Users{
			Active:     true,
			Email:      "frieda@oberbrunnerkirlin.name",
			ExternalID: externalID,
			ID:         "5975c461f9f8eb02aae053f3",
			Roles:      []string{"admin", "user"},
		}, nil
	}
	if email == "example@notexists.com" {
		return nil, goa.ErrNotFound("user not exist")
	}
	if email == "example@invalid.com" {
		return nil, goa.ErrInternal("internal server error")
	}

	for userID, user := range db.users {
		exID := ""
		if user.ExternalID != nil {
			exID = *user.ExternalID
		}
		if user.Email == email {
			return &app.Users{
				Active:        user.Active,
				Email:         user.Email,
				ExternalID:    exID,
				ID:            userID,
				Organizations: user.Organizations,
				Roles:         user.Roles,
			}, nil
		}
	}

	return nil, nil
}

// ActivateUser mock activation of a user.
func (db *DB) ActivateUser(email string) error {
	if email == "trigger-server-error@example.com" {
		return goa.ErrInternal("intentional server error")
	}
	var user *app.UserPayload
	for _, u := range db.users {
		if u.Email == email {
			user = u
			break
		}
	}
	if user == nil {
		return goa.ErrNotFound("not found")
	}
	user.Active = true
	return nil
}

// TokensMock implements a mock of ITokenCollection.
type TokensMock struct {
	Tokens map[string]*app.UserPayload
}

// CreateToken creates a token entry in the mock.
func (m *TokensMock) CreateToken(payload *app.UserPayload) error {
	m.Tokens[*payload.Token] = payload
	return nil
}

// VerifyToken performs a mock verification of a token.
func (m *TokensMock) VerifyToken(token string) (*string, error) {
	payload, ok := m.Tokens[token]
	if !ok {
		return nil, goa.ErrNotFound("not found")
	}
	return &payload.Email, nil
}

// DeleteToken removes a token record from the mock.
func (m *TokensMock) DeleteToken(token string) error {
	_, ok := m.Tokens[token]
	if !ok {
		return goa.ErrNotFound("not found")
	}
	delete(m.Tokens, token)
	return nil
}

// DeleteUserToken removes all tooken records from the mock for a user with the given email.
func (m *TokensMock) DeleteUserToken(email string) error {
	deleteTokens := []string{}
	for token, payload := range m.Tokens {
		if payload.Email == email {
			deleteTokens = append(deleteTokens, token)
		}
	}

	for _, token := range deleteTokens {
		delete(m.Tokens, token)
	}
	return nil
}
