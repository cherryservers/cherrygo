package cherrygo

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"time"
)

const lbPath = "/v1/load-balancers"

// LoadBalancerService is used to manage load balancer resources.
//
// Check out the [API docs] for specifics on the API
// and [product docs] for details on how Cherry Servers load balancers work.
//
// [API docs]: https://api.cherryservers.com/doc/#tag/Load-Balancer
// [product docs]: https://www.cherryservers.com/knowledge/docs/networking/load-balancer
type LoadBalancerService interface {
	Get(ctx context.Context, id int, opts *GetOptions) (LoadBalancer, *Response, error)
	List(ctx context.Context, projectID int, opts *GetOptions) ([]LoadBalancer, *Response, error)
	Create(ctx context.Context, projectID int, request CreateLoadBalancer) (LoadBalancer, *Response, error)
	Delete(ctx context.Context, id int) (*Response, error)
	ListPlans(ctx context.Context, teamID int, opts *GetOptions) ([]LoadBalancerPlan, *Response, error)
	Update(ctx context.Context, id int, request UpdateLoadBalancer) (LoadBalancer, *Response, error)
	Reset(ctx context.Context, id int) (LoadBalancer, *Response, error)
}

// LoadBalancer represents a Cherry Servers load balancer resource.
//
// Check out the [product docs] for details on how Cherry Servers load balancers work.
//
// [product docs]: https://www.cherryservers.com/knowledge/docs/networking/load-balancer
type LoadBalancer struct {
	ID     int              `json:"id,omitzero"`
	Name   string           `json:"name,omitzero"`
	Status string           `json:"status,omitzero"`
	Plan   LoadBalancerPlan `json:"plan,omitzero"`

	// Servers that the load balancer may route traffic to.
	Servers []Server `json:"servers,omitzero"`

	// Rules for routing traffic to backend servers.
	Rules      []LoadBalancerRule `json:"rules,omitzero"`
	IPAdresses []IPAddress        `json:"ip_addresses,omitzero"`
	Region     Region             `json:"region,omitzero"`

	// Backends define the rules and mechanisms for routing traffic
	// to a particular backend server.
	Backends []LoadBalancerBackend `json:"backends,omitzero"`

	// StickyCookie is the cookie that is used to route subsequent requests from
	// the same client to the same server.
	StickyCookie string `json:"sticky_cookie,omitzero"`

	// StickyEnabled indicates whether sticky sessions are enabled with the
	// use of StickyCookie.
	StickyEnabled bool `json:"sticky_enabled"`

	// HTTPSRedirectEnabled indicates whether all HTTP traffic will be
	// redirected through HTTPS via a 307 redirect.
	HTTPSRedirectEnabled bool `json:"https_redirect_enabled"`

	// ProxyEnabled indicates whether the client IP address will be preserved
	// as the request passes through the load balancer.
	ProxyEnabled bool `json:"proxy_enabled"`

	// Certificates that the load balancer uses, if HTTPS routing
	// rules are in effect.
	Certificates []LoadBalancerCertificate `json:"certificates,omitzero"`
	Pricing      Pricing                   `json:"pricing,omitzero"`

	// HealthCheck defines the health checking rules for backend servers.
	HealthCheck LoadBalancerHealthCheck `json:"health_check,omitzero"`
}

// LoadBalancerPlan represents the core load balancer attributes,
// based on price and regional availability.
type LoadBalancerPlan struct {
	Slug       string                      `json:"slug,omitzero"`
	Pricing    []Pricing                   `json:"pricing,omitzero"`
	Regions    []Region                    `json:"regions,omitzero"`
	Attributes []LoadBalancerPlanAttribute `json:"attributes,omitzero"`
}

// LoadBalancerPlanAttribute defines load balancer attributes,
// such as max RPS or bandwidth limitations.
type LoadBalancerPlanAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// LoadBalancerRule defines how traffic is routed to backends.
type LoadBalancerRule struct {
	ID              string `json:"id,omitzero"`
	SourcePort      int    `json:"source_port,omitzero"`
	DestinationPort int    `json:"destination_port,omitzero"`

	// SourceProtocol supports "tcp", "http", "http2" and "https".
	// "https" requires adding a certificate.
	SourceProtocol string `json:"source_protocol,omitzero"`

	// DestinationProtocol supports "tcp", "http", "http2" and "https".
	// "https" requires adding a certificate.
	DestinationProtocol string `json:"destination_protocol,omitzero"`

	// StickyEnabled indicates if sticky sessions (the use of a cookie
	// to route subsequent requests from a client to the same server) are enabled.
	StickyEnabled bool `json:"sticky_enabled"`

	// Balance defines the algorithm to use for traffic balancing.
	// The options are:
	//   - roundrobin - forwards requests to each backend according
	//     to their assigned weight; the higher the weight, the more requests are sent.
	//   - static-rr - forwards requests to each backend in sequential order.
	//   - leastconn - forwards requests to the backend that had the least active
	//     connections at the time the request was made.
	//   - source - preserves the incoming request’s original IP when forwarding it to the backend.
	Balance string `json:"balance,omitzero"`
	Status  string `json:"status,omitzero"`
}

// LoadBalancerBackend defines a backend server to which traffic is routed,
// as well as the rule that specifies the routing mechanism.
type LoadBalancerBackend struct {
	Rule LoadBalancerRule `json:"rule,omitzero"`
	IP   netip.Addr       `json:"ip,omitzero"`
	Port int              `json:"port,omitzero"`

	// Status should be "Up" for a healthy backend server.
	Status string `json:"state,omitzero"`

	// StatusDetail is a more detailed status description, such as
	// "Layer4 connection problem".
	StatusDetail string `json:"status,omitzero"`

	// Weigh is used by some of the balance options in Rule
	// to determine traffic distribution.
	Weight int `json:"weight,omitzero"`

	// Fails indicates how many times the backend server has failed to respond.
	Fails int `json:"fails,omitzero"`

	// Downs indicates how many times the backend servers went down.
	Downs int `json:"downs,omitzero"`

	// Downtime in seconds.
	Downtime int `json:"downtime,omitzero"`
}

// LoadBalancerCertificate is a certificate for HTTPS-based rules.
type LoadBalancerCertificate struct {
	ID      string    `json:"id,omitzero"`
	CN      string    `json:"cn,omitzero"`
	Starts  time.Time `json:"starts,omitzero"`
	Expires time.Time `json:"expires,omitzero"`
}

// LoadBalancerHealthCheck defines the parameters for health checks, should it be enabled.
type LoadBalancerHealthCheck struct {
	Enabled bool `json:"enabled"`

	// Path defines the path at which to perform the health check.
	// Only valid for HTTP/HTTPS based rules.
	Path string `json:"path,omitzero"`

	// Interval in seconds.
	Interval int `json:"interval,omitzero"`

	// HealthyThreshold is the number of successful health checks for the
	// backend server, before it's considered healthy.
	HealthyThreshold int `json:"healthy_threshold,omitzero"`

	// UnhealthyThreshold is the number of failed health checks for the
	// backend server, before it's considered unhealthy.
	UnhealthyThreshold int `json:"unhealthy_threshold,omitzero"`
}

// LoadBalancerClient makes load balancer related API requests.
type LoadBalancerClient struct {
	client *Client
}

// Get retrieves a load balancer resource.
func (c *LoadBalancerClient) Get(ctx context.Context, id int, opts *GetOptions) (LoadBalancer, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", lbPath, id))
	var lb LoadBalancer

	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return LoadBalancer{}, nil, err
	}

	resp, err := c.client.Do(req, &lb)
	return lb, resp, err
}

// List retrieves all load balancers in a project.
func (c *LoadBalancerClient) List(ctx context.Context, projectID int, opts *GetOptions) ([]LoadBalancer, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/load-balancers", baseProjectPath, projectID))
	var lbs []LoadBalancer

	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []LoadBalancer{}, nil, err
	}

	resp, err := c.client.Do(req, &lbs)
	return lbs, resp, err
}

// CreateLoadBalancer is the body for a load balancer creation request.
type CreateLoadBalancer struct {
	// Name for the load balancer resource. Will be auto-generated, if not set.
	Name string `json:"name,omitzero"`

	// Plan slug. Required.
	Plan string `json:"slug"`

	// Region slug. Note that not all regions may support load balancers,
	// see the [product docs]. Defaults to `LT-Siauliai`.
	//
	// [product docs]: https://www.cherryservers.com/knowledge/docs/networking/load-balancer
	Region string `json:"region,omitzero"`

	// Cycle sets the billing cycle for the load balancer.
	// Defaults to "hourly". If another cycle is chosen, funds will be charged
	// from the account balance to cover the invoice.
	Cycle string `json:"cycle,omitzero"`
}

// Create a new load balancer.
func (c *LoadBalancerClient) Create(ctx context.Context, projectID int, request CreateLoadBalancer) (LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%d/load-balancers", baseProjectPath, projectID)
	var lb LoadBalancer

	req, err := c.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return LoadBalancer{}, nil, err
	}

	resp, err := c.client.Do(req, &lb)
	return lb, resp, err
}

// Delete a load balancer.
func (c *LoadBalancerClient) Delete(ctx context.Context, id int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", lbPath, id)

	req, err := c.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req, nil)
}

// ListPlans lists load balancer plans.
func (c *LoadBalancerClient) ListPlans(ctx context.Context, teamID int, opts *GetOptions) ([]LoadBalancerPlan, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/load-balancer-plans", teamsPath, teamID))
	var plans []LoadBalancerPlan

	req, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return []LoadBalancerPlan{}, nil, err
	}

	resp, err := c.client.Do(req, &plans)
	return plans, resp, err
}

// UpdateLoadBalancer is the body for a load balancer update request.
type UpdateLoadBalancer struct {
	Name string `json:"name,omitzero"`

	// Plan slug. Initiates an upgrade process, during which the load balancer will be unavailable.
	Plan string `json:"slug,omitzero"`

	// StickyCookie used by sticky sessions.
	StickyCookie string `json:"sticky_cookie,omitzero"`

	// StickEnabled enables sticky sessions, which route subsequent requests
	// from the same client to the same server.
	StickyEnabled *bool `json:"sticky_enabled,omitzero"`

	// HTTPSRedirectEnabled redirects all HTTP traffic through HTTPS with a 307 redirect.
	HTTPSRedirectEnabled *bool `json:"https_redirect_enabled,omitzero"`

	// ProxyEnabled preserves client IP as the request passes through the load balancer.
	ProxyEnabled *bool `json:"proxy_enabled,omitzero"`

	// HealthCheckEnabled enables backend server health checks. Only for HTTP/S rules.
	HealthCheckEnabled *bool `json:"health_check_enabled,omitzero"`

	// HealthCheckPath is the path at which the health checks are performed.
	HealthCheckPath string `json:"health_check_path,omitzero"`

	// HealthCheckInterval determines how often to perform health checks, in seconds.
	HealthCheckInterval int `json:"health_check_interval,omitzero"`

	// HealthyThreshold determines how many successful health checks are required
	// for a backend server to become healthy.
	HealthyThreshold int `json:"healthy_threshold,omitzero"`

	// UnhealthyThreshold determines how many failed health checks are required for
	// a backend server to become unhealthy.
	UnhealthyThreshold int `json:"unhealthy_threshold,omitzero"`
}

// Update load balancer.
func (c *LoadBalancerClient) Update(ctx context.Context, id int, request UpdateLoadBalancer) (LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%d", lbPath, id)
	var lb LoadBalancer

	req, err := c.client.NewRequest(ctx, http.MethodPut, path, request)
	if err != nil {
		return LoadBalancer{}, nil, err
	}

	resp, err := c.client.Do(req, &lb)
	return lb, resp, err
}

// Reset load balancer.
// This will re-deploy the load balancer, keeping all existing configuration.
func (c *LoadBalancerClient) Reset(ctx context.Context, id int) (LoadBalancer, *Response, error) {
	path := fmt.Sprintf("%s/%d/actions", lbPath, id)
	bod := struct {
		Type string `json:"type"`
	}{
		Type: "reset",
	}
	var lb LoadBalancer

	req, err := c.client.NewRequest(ctx, http.MethodPost, path, bod)
	if err != nil {
		return LoadBalancer{}, nil, err
	}

	resp, err := c.client.Do(req, &lb)
	return lb, resp, err
}
