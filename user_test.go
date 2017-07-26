package main

import (
	"testing"
	"context"

	"user-microservice/app/test"
	"user-microservice/store"
	"github.com/goadesign/goa"
)

var (
	service = goa.New("user-test")
	db      = store.NewDB()
	ctrl    = NewUserController(service, db)
)

func TestGetUserOK(t *testing.T) {
	// Call generated test helper, this checks that the returned media type is of the
	// correct type (i.e. uses view "default") and validates the media type.
	// Also, it ckecks the returned status code
	_, user := test.GetUserOK(t, context.Background(), service, ctrl, "5975c461f9f8eb02aae053f3")
	
	if user == nil {
		t.Fatal("Nil user")
	}

	if user.ID != "5975c461f9f8eb02aae053f3" {
		t.Errorf("Invalid user ID, expected 5975c461f9f8eb02aae053f3, got %s", user.ID)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// The test helper takes care of validating the status code for us
	test.GetUserNotFound(t, context.Background(), service, ctrl, "fakeobjectidab02aae053f3")	
}