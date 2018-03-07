// Code generated by goagen v1.3.0, DO NOT EDIT.
//
// API "user": user Resource Client
//
// Command:
// $ goagen
// --design=github.com/JormungandrK/user-microservice/design
// --out=$(GOPATH)/src/github.com/JormungandrK/user-microservice
// --version=v1.3.0

package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// CreateUserPath computes a request path to the create action of user.
func CreateUserPath() string {

	return fmt.Sprintf("/users")
}

// Creates user
func (c *Client) CreateUser(ctx context.Context, path string, payload *UserPayload, contentType string) (*http.Response, error) {
	req, err := c.NewCreateUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewCreateUserRequest create the request corresponding to the create action endpoint of the user resource.
func (c *Client) NewCreateUserRequest(ctx context.Context, path string, payload *UserPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// FindUserPath computes a request path to the find action of user.
func FindUserPath() string {

	return fmt.Sprintf("/users/find")
}

// Find a user by email+password
func (c *Client) FindUser(ctx context.Context, path string, payload *Credentials, contentType string) (*http.Response, error) {
	req, err := c.NewFindUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewFindUserRequest create the request corresponding to the find action endpoint of the user resource.
func (c *Client) NewFindUserRequest(ctx context.Context, path string, payload *Credentials, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// FindByEmailUserPath computes a request path to the findByEmail action of user.
func FindByEmailUserPath() string {

	return fmt.Sprintf("/users/find/email")
}

// Find a user by email
func (c *Client) FindByEmailUser(ctx context.Context, path string, payload *EmailPayload, contentType string) (*http.Response, error) {
	req, err := c.NewFindByEmailUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewFindByEmailUserRequest create the request corresponding to the findByEmail action endpoint of the user resource.
func (c *Client) NewFindByEmailUserRequest(ctx context.Context, path string, payload *EmailPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// GetUserPath computes a request path to the get action of user.
func GetUserPath(userID string) string {
	param0 := userID

	return fmt.Sprintf("/users/%s", param0)
}

// Get user by id
func (c *Client) GetUser(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetUserRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetUserRequest create the request corresponding to the get action endpoint of the user resource.
func (c *Client) NewGetUserRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// GetMeUserPath computes a request path to the getMe action of user.
func GetMeUserPath() string {

	return fmt.Sprintf("/users/me")
}

// Retrieves the user information for the authenticated user
func (c *Client) GetMeUser(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewGetMeUserRequest(ctx, path)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewGetMeUserRequest create the request corresponding to the getMe action endpoint of the user resource.
func (c *Client) NewGetMeUserRequest(ctx context.Context, path string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// ResetVerificationTokenUserPath computes a request path to the resetVerificationToken action of user.
func ResetVerificationTokenUserPath() string {

	return fmt.Sprintf("/users/verification/reset")
}

// Reset verification token
func (c *Client) ResetVerificationTokenUser(ctx context.Context, path string, payload *EmailPayload, contentType string) (*http.Response, error) {
	req, err := c.NewResetVerificationTokenUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewResetVerificationTokenUserRequest create the request corresponding to the resetVerificationToken action endpoint of the user resource.
func (c *Client) NewResetVerificationTokenUserRequest(ctx context.Context, path string, payload *EmailPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("POST", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// UpdateUserPath computes a request path to the update action of user.
func UpdateUserPath(userID string) string {
	param0 := userID

	return fmt.Sprintf("/users/%s", param0)
}

// Update user
func (c *Client) UpdateUser(ctx context.Context, path string, payload *UserPayload, contentType string) (*http.Response, error) {
	req, err := c.NewUpdateUserRequest(ctx, path, payload, contentType)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewUpdateUserRequest create the request corresponding to the update action endpoint of the user resource.
func (c *Client) NewUpdateUserRequest(ctx context.Context, path string, payload *UserPayload, contentType string) (*http.Request, error) {
	var body bytes.Buffer
	if contentType == "" {
		contentType = "*/*" // Use default encoder
	}
	err := c.Encoder.Encode(payload, &body, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body: %s", err)
	}
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	req, err := http.NewRequest("PUT", u.String(), &body)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if contentType == "*/*" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", contentType)
	}
	return req, nil
}

// VerifyUserPath computes a request path to the verify action of user.
func VerifyUserPath() string {

	return fmt.Sprintf("/users/verify")
}

// Verify a user by token
func (c *Client) VerifyUser(ctx context.Context, path string, token *string) (*http.Response, error) {
	req, err := c.NewVerifyUserRequest(ctx, path, token)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(ctx, req)
}

// NewVerifyUserRequest create the request corresponding to the verify action endpoint of the user resource.
func (c *Client) NewVerifyUserRequest(ctx context.Context, path string, token *string) (*http.Request, error) {
	scheme := c.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: c.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if token != nil {
		values.Set("token", *token)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}
