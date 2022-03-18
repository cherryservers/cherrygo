package cherrygo

import (
	"fmt"
	"strconv"
	"strings"
)

const baseProjectPath = "/v1/teams"
const endProjectPath = "projects"

// GetProjects interface metodas isgauti team'sus
type GetProjects interface {
	List(teamID int, opts *GetOptions) ([]Project, *Response, error)
}

// ProjectsClient paveldi client
type ProjectsClient struct {
	client *Client
}

// List func lists teams
func (p *ProjectsClient) List(teamID int, opts *GetOptions) ([]Project, *Response, error) {
	//root := new(teamRoot)

	teamIDString := strconv.Itoa(teamID)

	path := strings.Join([]string{baseProjectPath, teamIDString, endProjectPath}, "/")
	pathQuery := opts.WithQuery(path)

	var trans []Project

	resp, err := p.client.MakeRequest("GET", pathQuery, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
