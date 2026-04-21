package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseUserPath = "/v1/users"

// UsersService is an interface for interfacing with the the User endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Users
type UsersService interface {
	CurrentUser(ctx context.Context, opts *GetOptions) (User, *Response, error)
	Get(ctx context.Context, userID int, opts *GetOptions) (User, *Response, error)
}

// User is the Cherry Servers user account.
type User struct {
	ID                    int    `json:"id,omitempty"`
	FirstName             string `json:"first_name,omitempty"`
	LastName              string `json:"last_name,omitempty"`
	Email                 string `json:"email,omitempty"`
	EmailVerified         bool   `json:"email_verified,omitempty"`
	Phone                 string `json:"phone,omitempty"`
	SecurityPhone         string `json:"security_phone,omitempty"`
	SecurityPhoneVerified bool   `json:"security_phone_verified,omitempty"`
	Href                  string `json:"href,omitempty"`
}

// UsersClient makes user related API requests.
type UsersClient struct {
	client *Client
}

// CurrentUser gets current user based on the bearer token.
func (s *UsersClient) CurrentUser(ctx context.Context, opts *GetOptions) (User, *Response, error) {
	var trans User
	path := opts.WithQuery("v1/user")

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return User{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Get a user.
func (s *UsersClient) Get(ctx context.Context, userID int, opts *GetOptions) (User, *Response, error) {
	var trans User
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseUserPath, userID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return User{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}
