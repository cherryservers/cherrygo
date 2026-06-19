package cherrygo

import (
	"context"
	"fmt"
	"net/http"
)

const baseProjectPath = "/v1/projects"

// ProjectsService is an interface for interfacing with the Projects endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Projects
type ProjectsService interface {
	List(ctx context.Context, teamID int, opts *GetOptions) ([]Project, *Response, error)
	Get(ctx context.Context, projectID int, opts *GetOptions) (Project, *Response, error)
	Create(ctx context.Context, teamID int, request *CreateProject) (Project, *Response, error)
	Update(ctx context.Context, projectID int, request *UpdateProject) (Project, *Response, error)
	ListSSHKeys(ctx context.Context, projectID int, opts *GetOptions) ([]SSHKey, *Response, error)
	Delete(ctx context.Context, projectID int) (*Response, error)
}

// Project data.
type Project struct {
	ID   int        `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`
	Bgp  ProjectBGP `json:"bgp,omitempty"`
	Href string     `json:"href,omitempty"`
}

// CreateProject fields for adding new project with specified name
type CreateProject struct {
	Name string `json:"name,omitempty"`
	Bgp  bool   `json:"bgp,omitempty"`
}

// UpdateProject fields for updating a project with specified name
type UpdateProject struct {
	Name *string `json:"name,omitempty"`
	Bgp  *bool   `json:"bgp,omitempty"`
}

// ProjectsClient makes project related API requests.
type ProjectsClient struct {
	client *Client
}

// List func lists projects
func (p *ProjectsClient) List(ctx context.Context, teamID int, opts *GetOptions) ([]Project, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/teams/%d/projects", teamID))
	var trans []Project

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}

// Get project.
func (p *ProjectsClient) Get(ctx context.Context, projectID int, opts *GetOptions) (Project, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseProjectPath, projectID))
	var trans Project

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Project{}, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}

// Create func will create new Project for specified team
func (p *ProjectsClient) Create(ctx context.Context, teamID int, request *CreateProject) (Project, *Response, error) {
	var trans Project
	path := fmt.Sprintf("/v1/teams/%d/projects", teamID)

	req, err := p.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return Project{}, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}

// Update func will update a project
func (p *ProjectsClient) Update(ctx context.Context, projectID int, request *UpdateProject) (Project, *Response, error) {
	var trans Project
	path := fmt.Sprintf("%s/%d", baseProjectPath, projectID)

	req, err := p.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return Project{}, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}

// Delete func will delete a project
func (p *ProjectsClient) Delete(ctx context.Context, projectID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", baseProjectPath, projectID)

	req, err := p.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req, nil)
	return resp, err
}

// ListSSHKeys available for project.
func (p *ProjectsClient) ListSSHKeys(ctx context.Context, projectID int, opts *GetOptions) ([]SSHKey, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/projects/%d/ssh-keys", projectID))
	var trans []SSHKey

	req, err := p.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := p.client.Do(req, &trans)
	return trans, resp, err
}
