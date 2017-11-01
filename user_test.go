package main

import (
	"context"
	"testing"

	"github.com/JormungandrK/microservice-security/auth"
	"github.com/JormungandrK/user-microservice/app"
	"github.com/JormungandrK/user-microservice/app/test"
	"github.com/JormungandrK/user-microservice/store"
	"github.com/goadesign/goa"
)

var (
	service          = goa.New("user-test")
	db               = store.NewDB()
	ctrl             = NewUserController(service, db)
	hexObjectID      = "5975c461f9f8eb02aae053f3"
	fakeHexObjectID  = "6975c461f9f8eb02aae053f3"
	badHexObjectID   = "fakeobjectidab02aae053f3aasadas"
	internalErrObjID = "6975c461f9f8eb02aae053f4"
)

func TestGetUserOK(t *testing.T) {
	// Call generated test helper, this checks that the returned media type is of the
	// correct type (i.e. uses view "default") and validates the media type.
	// Also, it ckecks the returned status code
	_, user := test.GetUserOK(t, context.Background(), service, ctrl, hexObjectID)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != hexObjectID {
		t.Errorf("Invalid user ID, expected %s, got %s", hexObjectID, user.ID)
	}
}

// The test helper takes care of validating the status code for us
func TestGetUserNotFound(t *testing.T) {
	test.GetUserNotFound(t, context.Background(), service, ctrl, fakeHexObjectID)
}

func TestGetUserBadRequest(t *testing.T) {
	test.GetUserBadRequest(t, context.Background(), service, ctrl, badHexObjectID)
}

func TestGetUserInternalServerError(t *testing.T) {
	test.GetUserInternalServerError(t, context.Background(), service, ctrl, internalErrObjID)
}

func TestGetMeUserOK(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: hexObjectID}
	ctx = auth.SetAuth(ctx, authObj)

	_, user := test.GetMeUserOK(t, ctx, service, ctrl)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != hexObjectID {
		t.Errorf("Invalid user ID: expected %s, got %s", hexObjectID, user.ID)
	}
}

func TestGetMeUserNotFound(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: fakeHexObjectID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetMeUserNotFound(t, ctx, service, ctrl)
}

func TestGetMeUserBadRequest(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: badHexObjectID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetMeUserBadRequest(t, ctx, service, ctrl)
}

func TestGetMeUserInternalServerError(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: internalErrObjID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetMeUserInternalServerError(t, ctx, service, ctrl)
}

func TestCreateUserOK(t *testing.T) {
	roles := []string{"admin", "user"}
	pass := "password"
	extID := "qwerc461f9f8eb02aae053f3"
	userPayload := &app.UserPayload{
		Password:   &pass,
		Email:      "example@some.com",
		ExternalID: &extID,
		Roles:      roles,
	}

	//CreateUserCreated
	_, user := test.CreateUserCreated(t, context.Background(), service, ctrl, userPayload)

	if user == nil {
		t.Fatal("User not created")
	}
}

func TestCreateUserBadRequest(t *testing.T) {
	userPayload := &app.UserPayload{
		Email: "example@some.com",
		Roles: []string{"admin", "user"},
	}

	test.CreateUserBadRequest(t, context.Background(), service, ctrl, userPayload)
}

func TestCreateUserInternalServerError(t *testing.T) {
	pass := "password"
	extID := "qwerc461f9f8eb02aae053f3"
	userPayload := &app.UserPayload{
		Password:   &pass,
		Email:      "internal-error@example.com",
		ExternalID: &extID,
		Roles:      []string{"admin", "user"},
	}

	test.CreateUserInternalServerError(t, context.Background(), service, ctrl, userPayload)
}

func TestUpdateUserOK(t *testing.T) {
	roles := []string{"admin", "user"}
	pass := "password"
	extID := "qwerc461f9f8eb02aae053f3"
	userPayload := &app.UserPayload{
		Password:   &pass,
		Email:      "example@some.com",
		ExternalID: &extID,
		Roles:      roles,
	}
	_, users := test.UpdateUserOK(t, context.Background(), service, ctrl, hexObjectID, userPayload)
	if users == nil {
		t.Fatal("Expected the update user data.")
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	pass := "password"
	extID := "qwerc461f9f8eb02aae053f3"
	userPayload := &app.UserPayload{
		Password:   &pass,
		Email:      "example@some.com",
		ExternalID: &extID,
		Roles:      []string{"admin", "user"},
	}

	test.UpdateUserNotFound(t, context.Background(), service, ctrl, fakeHexObjectID, userPayload)
}

func TestUpdateUserBadRequest(t *testing.T) {
	pass := "password"
	extID := "qwerc461f9f8eb02aae053f3"
	userPayload := &app.UserPayload{
		Password:   &pass,
		Email:      "example@some.com",
		ExternalID: &extID,
		Roles:      []string{"admin", "user"},
	}

	test.UpdateUserBadRequest(t, context.Background(), service, ctrl, badHexObjectID, userPayload)
}

func TestUpdateUserInternalServerError(t *testing.T) {
	pass := "password"
	extID := "qwerc461f9f8eb02aae053f3"
	userPayload := &app.UserPayload{
		Password:   &pass,
		Email:      "example@some.com",
		ExternalID: &extID,
		Roles:      []string{"admin", "user"},
	}

	test.UpdateUserInternalServerError(t, context.Background(), service, ctrl, internalErrObjID, userPayload)
}

func TestFindUserBadRequest(t *testing.T) {
	payload := &app.Credentials{
		Email:    "",
		Password: "",
	}
	test.FindUserBadRequest(t, context.Background(), service, ctrl, payload)
}

func TestFindUserInternalServerError(t *testing.T) {
	payload := &app.Credentials{
		Email:    "internal-error@example.com",
		Password: "the-pass",
	}
	test.FindUserInternalServerError(t, context.Background(), service, ctrl, payload)
}

func TestFindUserNotFound(t *testing.T) {
	payload := &app.Credentials{
		Email:    "example@notexists.com",
		Password: "the-pass",
	}
	test.FindUserNotFound(t, context.Background(), service, ctrl, payload)
}

func TestFindUserOK(t *testing.T) {
	payload := &app.Credentials{
		Email:    "email@example.com",
		Password: "valid-pass",
	}
	_, user := test.FindUserOK(t, context.Background(), service, ctrl, payload)
	if user == nil {
		t.Fatal("Expected user")
	}
}

func TestFindByEmailUserOK(t *testing.T) {
	payload := &app.EmailPayload{
		Email: "frieda@oberbrunnerkirlin.name",
	}
	_, user := test.FindByEmailUserOK(t, context.Background(), service, ctrl, payload)

	if user == nil {
		t.Fatal("Nil user")
	}
}

func TestFindByEmailUserNotFound(t *testing.T) {
	payload := &app.EmailPayload{
		Email: "example@notexists.com",
	}

	test.FindByEmailUserNotFound(t, context.Background(), service, ctrl, payload)
}

func TestFindByEmailUserInternalServerError(t *testing.T) {
	payload := &app.EmailPayload{
		Email: "example@invalid.com",
	}

	test.FindByEmailUserInternalServerError(t, context.Background(), service, ctrl, payload)
}
