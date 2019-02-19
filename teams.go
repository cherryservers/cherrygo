package cherrygo

import (
	"log"
)

const teamsPath = "/v1/teams"

// GetTeams interface metodas isgauti team'sus
type GetTeams interface {
	List() ([]Teams, *Response, error)
}

// Kad teamsus i slice'us irasytume
type teamRoot struct {
	Teams []Teams
}

// Teams tai ka grazina api
type Teams struct {
	ID      int     `json:"id", omitempty`
	Name    string  `json:"name", omitempty`
	Credit  Credit  `json:"credit"`
	Billing Billing `json:"billing"`
	Href    string  `json:"href, omitempty`
}

// Credit fields
type Credit struct {
	Account   Account   `json:"account", ommitempty`
	Promo     Promo     `json:"promo", ommitempty`
	Resources Resources `json:"resources", ommitempty`
}

// Account fields
type Account struct {
	Remaining float32 `json:"remaining"`
	Usage     float32 `json:"usage"`
	Currency  string  `json:"currency"`
}

// Promo fields
type Promo struct {
	Remaining float32 `json:"remaining"`
	Usage     float32 `json:"usage"`
	Currency  string  `json:"currency"`
}

// Resources fields
type Resources struct {
	Pricing Pricing `json:"pricing", ommitempty`
}

// Pricing for resources
type Pricing struct {
	Price    float32 `json:"price", omitempty`
	Taxed    bool    `json:"taxed", omitempty`
	Currency string  `json:"currency", omitempty`
	Unit     string  `json:"unit", omitempty`
}

// Billing fields
type Billing struct {
	Type        string `json:"type", omitempty`
	CompanyName string `json:"company_name", ommitempty`
	CompanyCode string `json:"company_code", ommitempty`
	FirstName   string `json:"first_name", omitempty`
	LastName    string `json:"last_name", ommitempty`
	Address1    string `json:"address_1", ommitempty`
	Address2    string `json:"address_2", ommitempty`
	CountryIso2 string `json:"country_iso_2", ommitempty`
	City        string `json:"city", ommitempty`
	Vat         Vat    `json:"vat", ommitempty`
	Currency    string `json:"currency", ommitempty`
}

// Vat fields
type Vat struct {
	Amount int    `json:"amount"`
	Number string `json:"number, omitempty`
	Valid  bool   `json:"valid"`
}

// TeamsClient paveldi client
type TeamsClient struct {
	client *Client
}

// List func lists teams
func (t *TeamsClient) List() ([]Teams, *Response, error) {
	//root := new(teamRoot)

	var trans []Teams
	//resp := t.client.Bumba()
	//log.Println("\nFROM LIST1: ", root.Teams)
	resp, err := t.client.MakeRequest("GET", teamsPath, nil, &trans)
	if err != nil {
		log.Fatal(err)
	}

	return trans, resp, err
}
