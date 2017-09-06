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
	HexObjectID     = "5975c461f9f8eb02aae053f3"
	DiffHexObjectID = "6975c461f9f8eb02aae053f3"
	FakeHexObjectID = "fakeobjectidab02aae053f3"
)

func TestGetUserOK(t *testing.T) {
	// Call generated test helper, this checks that the returned media type is of the
	// correct type (i.e. uses view "default") and validates the media type.
	// Also, it ckecks the returned status code
	_, user := test.GetUserOK(t, context.Background(), service, ctrl, HexObjectID)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != HexObjectID {
		t.Errorf("Invalid user ID, expected %s, got %s", HexObjectID, user.ID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// The test helper takes care of validating the status code for us
	test.GetUserNotFound(t, context.Background(), service, ctrl, FakeHexObjectID)
}

func TestGetMeUserOK(t *testing.T) {
	_, user := test.GetMeUserOK(t, context.Background(), service, ctrl, &HexObjectID)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != HexObjectID {
		t.Errorf("Invalid user ID, expected %s, got %s", HexObjectID, user.ID)
	}

}

func TestGetMeUserNotFound(t *testing.T) {
	test.GetMeUserNotFound(t, context.Background(), service, ctrl, &FakeHexObjectID)
}

func TestCreateUserOK(t *testing.T) {
	roles := []string{"admin", "user"}
	userPayload := &app.UserPayload{
		Username:   "username",
		Password:   "password",
		Email:      "example@some.com",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		Roles:      roles,
	}

	//CreateUserCreated
	_, user := test.CreateUserCreated(t, context.Background(), service, ctrl, userPayload)

	if user == nil {
		t.Fatal("User not created")
	}
}

func TestUpdateUserOK(t *testing.T) {
	roles := []string{"admin", "user"}
	userPayload := &app.UserPayload{
		Username:   "username",
		Password:   "password",
		Email:      "example@some.com",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		Roles:      roles,
	}
	_, users := test.UpdateUserOK(t, context.Background(), service, ctrl, HexObjectID, userPayload)
	if users == nil {
		t.Fatal("Expected the update user data.")
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	roles := []string{"admin", "user"}
	userPayload := &app.UserPayload{
		Username:   "username",
		Password:   "password",
		Email:      "example@some.com",
		ExternalID: "qwerc461f9f8eb02aae053f3",
		Roles:      roles,
	}

	test.UpdateUserNotFound(t, context.Background(), service, ctrl, DiffHexObjectID, userPayload)
}

func TestFindUserBadRequest(t *testing.T) {
	payload := &app.Credentials{
		Username: "",
		Password: "",
	}
	test.FindUserBadRequest(t, context.Background(), service, ctrl, payload)
}

func TestFindUserInternalServerError(t *testing.T) {
	payload := &app.Credentials{
		Username: "internal-error-user",
		Password: "the-pass",
	}
	test.FindUserInternalServerError(t, context.Background(), service, ctrl, payload)
}

func TestFindUserNotFound(t *testing.T) {
	payload := &app.Credentials{
		Username: "nonexisting",
		Password: "the-pass",
	}
	test.FindUserNotFound(t, context.Background(), service, ctrl, payload)
}

func TestFindUserOK(t *testing.T) {
	payload := &app.Credentials{
		Username: "validuser",
		Password: "valid-pass",
	}
	_, user := test.FindUserOK(t, context.Background(), service, ctrl, payload)
	if user == nil {
		t.Fatal("Expected user")
	}
}
