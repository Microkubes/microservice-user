package main

import (
	"context"
	"testing"

	"github.com/JormungandrK/user-microservice/app"
	"github.com/JormungandrK/user-microservice/app/test"
	"github.com/JormungandrK/user-microservice/store"
	"github.com/goadesign/goa"
)

var (
	service         = goa.New("user-test")
	db              = store.NewDB()
	ctrl            = NewUserController(service, db)
	HexObjectId     = "5975c461f9f8eb02aae053f3"
	FakeHexObjectId = "fakeobjectidab02aae053f3"
)

func TestGetUserOK(t *testing.T) {
	// Call generated test helper, this checks that the returned media type is of the
	// correct type (i.e. uses view "default") and validates the media type.
	// Also, it ckecks the returned status code
	_, user := test.GetUserOK(t, context.Background(), service, ctrl, HexObjectId)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != HexObjectId {
		t.Errorf("Invalid user ID, expected %s, got %s", HexObjectId, user.ID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// The test helper takes care of validating the status code for us
	test.GetUserNotFound(t, context.Background(), service, ctrl, FakeHexObjectId)
}

func TestGetMeUserOK(t *testing.T) {
	_, user := test.GetMeUserOK(t, context.Background(), service, ctrl, &HexObjectId)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != HexObjectId {
		t.Errorf("Invalid user ID, expected %s, got %s", HexObjectId, user.ID)
	}

}

func TestGetMeUserNotFound(t *testing.T) {
	test.GetMeUserNotFound(t, context.Background(), service, ctrl, &FakeHexObjectId)
}

func TestCreateUserOK(t *testing.T) {
	roles := []string{"admin", "user"}
	userPayload := &app.UserPayload{
		Username:   "username",
		Password:   "password",
		Email:      "examlple@some.com",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		Roles:      roles,
	}

	//CreateUserCreated
	_, user := test.CreateUserCreated(t, context.Background(), service, ctrl, userPayload)

	if user == nil {
		t.Fatal("User not created")
	}
}
