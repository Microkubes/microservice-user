package main

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/Microkubes/microservice-security/auth"
	"github.com/Microkubes/microservice-user/app"
	"github.com/Microkubes/microservice-user/app/test"
	"github.com/Microkubes/microservice-user/store"
	"github.com/keitaroinc/goa"
)

var db = store.NewDB()
var (
	service          = goa.New("user-test")
	ctrl             = NewUserController(service, db, nil)
	ID               = "5df2103b5f1b640001142d3c"
	notFoundID       = "5df2103b5f1b640001142d4c"
	notFonundEmail   = "not-found@gmail.com"
	notFoundToken    = "not-found-token"
	badID            = "bad-id"
	badEmail         = "bad@gmail.com"
	internalErrID    = "internal-err-id"
	internalErrEmail = "internal-error@example.com"
	internalErrToken = "internal-error-token"
)

func TestGetUserOK(t *testing.T) {
	// Call generated test helper, this checks that the returned media type is of the
	// correct type (i.e. uses view "default") and validates the media type.
	// Also, it ckecks the returned status code
	_, user := test.GetUserOK(t, context.Background(), service, ctrl, ID)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != ID {
		t.Errorf("Invalid user ID, expected %s, got %s", ID, user.ID)
	}
}

// The test helper takes care of validating the status code for us
func TestGetUserNotFound(t *testing.T) {
	test.GetUserNotFound(t, context.Background(), service, ctrl, notFoundID)
}

func TestGetUserBadRequest(t *testing.T) {
	test.GetUserBadRequest(t, context.Background(), service, ctrl, badID)
}

func TestGetUserInternalServerError(t *testing.T) {
	test.GetUserInternalServerError(t, context.Background(), service, ctrl, internalErrID)
}

func TestGetMeUserOK(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: ID}
	ctx = auth.SetAuth(ctx, authObj)

	_, user := test.GetMeUserOK(t, ctx, service, ctrl)

	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != ID {
		t.Errorf("Invalid user ID: expected %s, got %s", ID, user.ID)
	}
}

func TestGetMeUserNotFound(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: notFoundID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetMeUserNotFound(t, ctx, service, ctrl)
}

func TestGetMeUserBadRequest(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: badID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetMeUserBadRequest(t, ctx, service, ctrl)
}

func TestGetMeUserInternalServerError(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: internalErrID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetMeUserInternalServerError(t, ctx, service, ctrl)
}

func TestCreateUserOK(t *testing.T) {
	roles := []string{"admin", "user"}
	password := "keitaro"
	extID := "qwerc461f9f8eb02aae053f3"
	CreateUserPayload := &app.CreateUserPayload{
		Email:      "keitaro-user2@gmail.com",
		Password:   &password,
		ExternalID: &extID,
		Roles:      roles,
	}

	//CreateUserCreated
	_, user := test.CreateUserCreated(t, context.Background(), service, ctrl, CreateUserPayload)

	if user == nil {
		t.Fatal("User not created")
	}
}

func TestCreateUserBadRequest(t *testing.T) {
	CreateUserPayload := &app.CreateUserPayload{
		Email: "example@some.com",
		Roles: []string{"admin", "user"},
	}

	test.CreateUserBadRequest(t, context.Background(), service, ctrl, CreateUserPayload)
}

func TestCreateUserInternalServerError(t *testing.T) {
	password := "keitaro"
	extID := "qwerc461f9f8eb02aae053f3"
	CreateUserPayload := &app.CreateUserPayload{
		Password:   &password,
		Email:      "internal-error@example.com",
		ExternalID: &extID,
		Roles:      []string{"admin", "user"},
	}

	test.CreateUserInternalServerError(t, context.Background(), service, ctrl, CreateUserPayload)
}

func TestUpdateUserOK(t *testing.T) {
	roles := []string{"admin"}
	UpdateUserPayload := &app.UpdateUserPayload{
		Roles: roles,
	}
	_, users := test.UpdateUserOK(t, context.Background(), service, ctrl, ID, UpdateUserPayload)
	if users == nil {
		t.Fatal("Expected the update user data.")
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	UpdateUserPayload := &app.UpdateUserPayload{
		Roles: []string{"admin", "user"},
	}

	test.UpdateUserNotFound(t, context.Background(), service, ctrl, notFoundID, UpdateUserPayload)
}

func TestUpdateUserBadRequest(t *testing.T) {
	UpdateUserPayload := &app.UpdateUserPayload{
		Roles: []string{"admin", "user"},
	}

	test.UpdateUserBadRequest(t, context.Background(), service, ctrl, badID, UpdateUserPayload)
}

func TestUpdateUserInternalServerError(t *testing.T) {
	UpdateUserPayload := &app.UpdateUserPayload{
		Roles: []string{"admin", "user"},
	}

	test.UpdateUserInternalServerError(t, context.Background(), service, ctrl, internalErrID, UpdateUserPayload)
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
		Email:    internalErrEmail,
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
		Email:    "keitaro-user2@gmail.com",
		Password: "keitaro",
	}
	_, user := test.FindUserOK(t, context.Background(), service, ctrl, payload)
	if user == nil {
		t.Fatal("Expected user")
	}
}

func TestFindByEmailUserOK(t *testing.T) {
	payload := &app.EmailPayload{
		Email: "keitaro-user1@gmail.com",
	}
	_, user := test.FindByEmailUserOK(t, context.Background(), service, ctrl, payload)

	if user == nil {
		t.Fatal("Nil user")
	}
}

func TestFindByEmailUserNotFound(t *testing.T) {
	payload := &app.EmailPayload{
		Email: notFonundEmail,
	}

	test.FindByEmailUserNotFound(t, context.Background(), service, ctrl, payload)
}

func TestFindByEmailUserInternalServerError(t *testing.T) {
	payload := &app.EmailPayload{
		Email: internalErrEmail,
	}

	test.FindByEmailUserInternalServerError(t, context.Background(), service, ctrl, payload)
}

func TestResetVerificationTokenUserOK(t *testing.T) {
	test.ResetVerificationTokenUserOK(t, context.Background(), service, ctrl, &app.EmailPayload{
		Email: "keitaro-user2@gmail.com",
	})
}

func TestResetVerificationTokenUserNotFound(t *testing.T) {
	test.ResetVerificationTokenUserNotFound(t, context.Background(), service, ctrl, &app.EmailPayload{
		Email: notFonundEmail,
	})
}

func TestResetVerificationTokenUserBadRequest(t *testing.T) {
	test.ResetVerificationTokenUserBadRequest(t, context.Background(), service, ctrl, &app.EmailPayload{
		Email: badEmail,
	})
	test.ResetVerificationTokenUserBadRequest(t, context.Background(), service, ctrl, &app.EmailPayload{})
}

func TestResetVerificationTokenUserInternalServerError(t *testing.T) {
	test.ResetVerificationTokenUserInternalServerError(t, context.Background(), service, ctrl, &app.EmailPayload{
		Email: internalErrEmail,
	})
}

func TestVerifyUserOK(t *testing.T) {
	token := "sdaewefdc234erfdd123erfdxc23edx"

	test.VerifyUserOK(t, context.Background(), service, ctrl, &token)
}

func TestVerifyUserNotFound(t *testing.T) {
	token := notFoundToken

	test.VerifyUserNotFound(t, context.Background(), service, ctrl, &token)
}

func TestVerifyUserInternalServerError(t *testing.T) {
	token := internalErrToken
	test.VerifyUserInternalServerError(t, context.Background(), service, ctrl, &token)
}

func TestGenerateToken(t *testing.T) {
	token := generateToken(40)

	if len(token) != 56 {
		t.Errorf("Expected token length was 56, got %d", len(token))
	}
}

func TestGenerateExpDate(t *testing.T) {
	expDate := generateExpDate()
	time := strconv.Itoa(int(time.Now().UTC().Unix()/60) + (60 * 24))
	if len(expDate) != len(time) {
		t.Errorf("expDate wrong format")
	}
	if expDate != time {
		t.Error("expDate value is not expected value")
	}
}

func TestCheckExpDate(t *testing.T) {
	expDate := strconv.Itoa(int(time.Now().UTC().Unix()/60) + (60 * 24))
	if !checkExpDate(expDate) {
		t.Errorf("expDate is expired, Expected value: [true]")
	}
	expDate = strconv.Itoa(int(time.Now().UTC().Unix()/60) - (60 * 24))
	if checkExpDate(expDate) {
		t.Errorf("expDate is not expired, Expected value [false]")
	}
}

func TestStringToBcryptHash(t *testing.T) {
	_, err := stringToBcryptHash("keitaro")

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestGetAllUserOK(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: ID}
	ctx = auth.SetAuth(ctx, authObj)

	test.GetAllUserOK(t, ctx, service, ctrl, nil, nil, nil, nil)
}

func TestGetAllUserNotFound(t *testing.T) {
	ctx := context.Background()
	authObj := &auth.Auth{UserID: ID}
	ctx = auth.SetAuth(ctx, authObj)

	offset := 5
	test.GetAllUserNotFound(t, ctx, service, ctrl, nil, &offset, nil, nil)
}

func TestGetAllUserInternalServerError(t *testing.T) {
	test.GetAllUserInternalServerError(t, context.Background(), service, ctrl, nil, nil, nil, nil)
}
