package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseSSHPath = "/v1/ssh-keys"

// SSHKeysService is an interface for interfacing with the the SSH keys endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/SshKeys
type SSHKeysService interface {
	List(ctx context.Context, opts *GetOptions) ([]SSHKey, *Response, error)
	Get(ctx context.Context, sshKeyID int, opts *GetOptions) (SSHKey, *Response, error)
	Create(ctx context.Context, request *CreateSSHKey) (SSHKey, *Response, error)
	Delete(ctx context.Context, sshKeyID int) (SSHKey, *Response, error)
	Update(ctx context.Context, sshKeyID int, request *UpdateSSHKey) (SSHKey, *Response, error)
}

// SSHKeys fields for return values after creation
type SSHKey struct {
	ID          int    `json:"id,omitempty"`
	Label       string `json:"label,omitempty"`
	Key         string `json:"key,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	User        User   `json:"user,omitempty"`
	Updated     string `json:"updated,omitempty"`
	Created     string `json:"created,omitempty"`
	Href        string `json:"href,omitempty"`
}

type SSHKeysClient struct {
	client *Client
}

// CreateSSHKey fields for adding new key with label and raw key
type CreateSSHKey struct {
	Label string `json:"label"`
	Key   string `json:"key"`
}

// UpdateSSHKey fields for label or key update
type UpdateSSHKey struct {
	Label *string `json:"label,omitempty"`
	Key   *string `json:"key,omitempty"`
}

// List all available ssh keys
func (s *SSHKeysClient) List(ctx context.Context, opts *GetOptions) ([]SSHKey, *Response, error) {
	var trans []SSHKey
	pathQuery := opts.WithQuery(baseSSHPath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, pathQuery, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

func (s *SSHKeysClient) Get(ctx context.Context, sshKeyID int, opts *GetOptions) (SSHKey, *Response, error) {
	var trans SSHKey
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseSSHPath, sshKeyID))

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return SSHKey{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Create adds new SSH key
func (s *SSHKeysClient) Create(ctx context.Context, request *CreateSSHKey) (SSHKey, *Response, error) {
	var trans SSHKey

	req, err := s.client.NewRequest(ctx, http.MethodPost, baseSSHPath, request)
	if err != nil {
		return SSHKey{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Delete removes desired SSH key by its ID
func (s *SSHKeysClient) Delete(ctx context.Context, sshKeyID int) (SSHKey, *Response, error) {
	var trans SSHKey
	path := fmt.Sprintf("%s/%d", baseSSHPath, sshKeyID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return SSHKey{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}

// Update function updates keys Label or key itself
func (s *SSHKeysClient) Update(ctx context.Context, sshKeyID int, request *UpdateSSHKey) (SSHKey, *Response, error) {
	var trans SSHKey
	path := fmt.Sprintf("%s/%d", baseSSHPath, sshKeyID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return SSHKey{}, nil, err
	}

	resp, err := s.client.Do(req, &trans)
	return trans, resp, err
}
