// Package cherrygo provides a client that can
// manage Cherry Servers infrastructure resources.
package cherrygo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/cherryservers/cherrygo/v3/internal/client"
)

const (
	apiURL             = "https://api.cherryservers.com/v1/"
	cherryAuthTokenVar = "CHERRY_AUTH_TOKEN"
	mediaType          = "application/json"
	userAgent          = "cherry-agent-go/"
	cherryDebugVar     = "CHERRY_DEBUG"
)

// Client is a client for the Cherry Servers RESTful API.
//
// Retries failed requests when it's safe to do so, e.g. status 429
// or a network timeout with an idempotent method. Respects `Retry-After` headers
// with fallback to exponential backoff with jitter.
type Client struct {
	client *client.Client

	BaseURL *url.URL

	UserAgent string
	AuthToken string

	Teams       TeamsService
	Plans       PlansService
	Images      ImagesService
	Projects    ProjectsService
	SSHKeys     SSHKeysService
	Servers     ServersService
	IPAddresses IPAddressesService
	Storages    StoragesService
	Regions     RegionsService
	Users       UsersService
	Backups     BackupsService
}

// Response is the http response from api calls.
type Response struct {
	*http.Response
	Meta
}

// Meta is the response metadata.
type Meta struct {
	Total int
}

// NewRequest creates a request. Adds the required headers.
func (c *Client) NewRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	url, _ := url.Parse(path)
	u := c.BaseURL.ResolveReference(url)

	buf := new(bytes.Buffer)
	if body != nil {
		coder := json.NewEncoder(buf)
		err := coder.Encode(body)
		if err != nil {
			log.Printf("Error while encoding body: %v -> %v", err, err.Error())
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + c.AuthToken
	req.Header.Set("Authorization", bearer)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", mediaType)
	if body != nil {
		req.Header.Add("Content-Type", mediaType)
	}
	return req, nil
}

// Do executes a request.
//
// The response body is un-marshalled into v, so it must be a pointer
// to a type that can hold the expected response, [io.Writer] or nil.
func (c *Client) Do(req *http.Request, v any) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := Response{Response: resp}
	response.populateTotal()

	if sc := response.StatusCode; sc >= 299 {
		type ErrorResponse struct {
			Response *http.Response
			Code     int    `json:"code"`
			Message  string `json:"message"`
		}

		var errorResponse ErrorResponse

		bod, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bod, &errorResponse)
		if err != nil {
			return nil, err
		}

		err = fmt.Errorf("error response from API: %v (error code: %v)", errorResponse.Message, errorResponse.Code)

		return &response, err
	}

	// Handling delete requests which EOF is not an error
	if req.Method == http.MethodDelete && response.StatusCode == http.StatusNoContent {
		return &response, err
	}

	if v != nil {
		// if v implements the io.Writer interface, return the raw response
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {

			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&v)
			if err != nil {
				log.Printf("error while decoding body: %v -> %v", err, err.Error())
				return &response, err
			}
		}
	}

	return &response, nil
}

type options struct {
	url       string
	client    *http.Client
	userAgent string
	authToken string
	debugDst  io.Writer
}

// ClientOpt is a client configuration option.
type ClientOpt func(*options) error

// NewClient creates a Cherry Servers API client.
func NewClient(opts ...ClientOpt) (*Client, error) {
	parsedOpts := &options{
		authToken: os.Getenv(cherryAuthTokenVar),
		client:    &http.Client{},
		url:       apiURL,
		userAgent: userAgent,
	}
	for _, opt := range opts {
		if err := opt(parsedOpts); err != nil {
			return nil, err
		}
	}
	if parsedOpts.authToken == "" {
		return nil, fmt.Errorf("auth token must be provided as parameter of environment variable %s", cherryAuthTokenVar)
	}

	url, err := url.Parse(parsedOpts.url)
	if err != nil {
		return nil, err
	}

	if parsedOpts.debugDst == nil && os.Getenv(cherryDebugVar) != "" {
		parsedOpts.debugDst = os.Stderr
	}

	c := &Client{
		client:    client.New(client.WithHTTPClient(parsedOpts.client), client.WithDebug(parsedOpts.debugDst)),
		AuthToken: parsedOpts.authToken,
		BaseURL:   url,
		UserAgent: parsedOpts.userAgent,
	}

	c.Teams = &TeamsClient{client: c}
	c.Plans = &PlansClient{client: c}
	c.Images = &ImagesClient{client: c}
	c.Projects = &ProjectsClient{client: c}
	c.SSHKeys = &SSHKeysClient{client: c}
	c.Servers = &ServersClient{client: c}
	c.IPAddresses = &IPsClient{client: c}
	c.Storages = &StoragesClient{client: c}
	c.Regions = &RegionsClient{client: c}
	c.Users = &UsersClient{client: c}
	c.Backups = &BackupsClient{client: c}

	return c, err
}

// ErrorResponse fields
type ErrorResponse struct {
	Response    *http.Response
	Errors      []string `json:"errors"`
	SingleError string   `json:"error"`
}

// WithUserAgent set user agent when making requests
func WithUserAgent(ua string) ClientOpt {
	return func(c *options) error {
		c.userAgent = fmt.Sprintf("%s %s", ua, userAgent)
		return nil
	}
}

// WithURL use url as endpoint for API requests
func WithURL(url string) ClientOpt {
	return func(c *options) error {
		c.url = url
		return nil
	}
}

// WithHTTPClient use client as the http.Client to make API requests
func WithHTTPClient(client *http.Client) ClientOpt {
	return func(c *options) error {
		c.client = client
		return nil
	}
}

// WithAuthToken use provided auth token to make requests, defaults to environment variable
// CHERRY_AUTH_TOKEN
func WithAuthToken(authToken string) ClientOpt {
	return func(c *options) error {
		c.authToken = authToken
		return nil
	}
}

// WithDebug enables debug mode, which dumps logs to w.
// Can also be enabled with the CHERRY_DEBUG environment variable, in which case logs
// will be dumped to stderr.
func WithDebug(w io.Writer) ClientOpt {
	return func(c *options) error {
		c.debugDst = w
		return nil
	}
}

func (r *Response) populateTotal() {
	// parse the headers and populate Meta.Total
	if total := r.Header.Get("X-Total-Count"); total != "" {
		r.Total, _ = strconv.Atoi(total)
	}
}
