package cherrygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTeam_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/teams/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodDelete)
		writer.WriteHeader(http.StatusNoContent)
		_, err := fmt.Fprint(writer)
		require.NoError(t, err)
	})

	_, err := testClient.Teams.Delete(t.Context(), 123)
	assert.NoError(t, err)
}

func TestTeam_List(t *testing.T) {
	setup()
	defer teardown()

	want := []Team{
		{
			ID:   123,
			Name: "team",
			Credit: Credit{
				Account: CreditDetails{
					Currency: "EUR",
				},
				Promo: CreditDetails{
					Remaining: 969.69,
					Usage:     189.72,
					Currency:  "EUR",
				},
			},
		},
	}

	mux.HandleFunc("/v1/teams", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		writer.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(writer, `[{
				"id":123,
				"name": "team",
				"credit": {
      				"account": {
       	 				"currency": "EUR"
      				},
					"promo": {
						"remaining": 969.69,
						"usage": 189.72,
						"currency": "EUR"
					}
				}
			 }]`)

		require.NoError(t, err)
	})

	team, _, err := testClient.Teams.List(t.Context(), nil)
	require.NoError(t, err)

	assert.Equal(t, want, team)
}

func TestTeam_Get(t *testing.T) {
	setup()
	defer teardown()

	want := Team{
		ID:   123,
		Name: "team",
		Credit: Credit{
			Account: CreditDetails{
				Currency: "EUR",
			},
			Promo: CreditDetails{
				Remaining: 969.48,
				Usage:     189.93,
				Currency:  "EUR",
			},
			Resources: Resources{
				Pricing: Pricing{
					Price:     0.2853,
					UnitPrice: 0,
					Taxed:     true,
					Currency:  "EUR",
					Unit:      "Hourly",
				},
				Remaining: RemainingTime{
					Time: 4107,
					Unit: "Hourly",
				},
			},
		},
		Billing: Billing{
			Type:        "personal",
			LastName:    "['test.test', 'cherryservers.com']",
			CountryIso2: "LT",
			Vat: Vat{
				Amount: 21,
				Valid:  false,
			},
			Currency: "EUR",
		},
		Href: "/teams/148226",
	}

	mux.HandleFunc("/v1/teams/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		writer.WriteHeader(http.StatusOK)

		_, err := fmt.Fprint(writer, `{
					"id": 123,
					"name": "team",
					"credit": {
						"account": {
						"currency": "EUR"
						},
						"promo": {
						"remaining": 969.48,
						"usage": 189.93,
						"currency": "EUR"
						},
						"resources": {
						"pricing": {
							"price": 0.2853,
							"unit_price": 0,
							"taxed": true,
							"currency": "EUR",
							"unit": "Hourly"
						},
						"remaining": {
							"time": 4107,
							"unit": "Hourly"
						}
						}
					},
					"billing": {
						"type": "personal",
						"last_name": "['test.test', 'cherryservers.com']",
						"country_iso_2": "LT",
						"vat": {
						"amount": 21,
						"valid": false
						},
						"currency": "EUR"
					},
					"href": "/teams/148226"
					}
		`)
		require.NoError(t, err)
	})

	team, _, err := testClient.Teams.Get(t.Context(), 123, nil)

	require.NoError(t, err)

	assert.Equal(t, want, team)
}

func TestTeam_Create(t *testing.T) {
	setup()
	defer teardown()

	want := Team{
		ID:   123,
		Name: "team",
	}

	createRequest := CreateTeam{
		Name:     "team",
		Type:     "personal",
		Currency: "EUR",
	}

	mux.HandleFunc("POST /v1/teams", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPost)
		v := new(CreateTeam)
		err := json.NewDecoder(request.Body).Decode(v)
		require.NoError(t, err)

		assert.Equal(t, createRequest, *v)

		writer.WriteHeader(http.StatusCreated)

		_, err = fmt.Fprint(writer, `{
					"id": 123,
					"name": "team"
			}
		`)
		require.NoError(t, err)
	})

	team, _, err := testClient.Teams.Create(t.Context(), &createRequest)
	require.NoError(t, err)

	assert.Equal(t, want, team)
}

func TestTeam_Update(t *testing.T) {
	setup()
	defer teardown()

	want := Team{
		ID:   123,
		Name: "team",
	}

	var (
		name     = "team"
		teamType = "personal"
		currency = "EUR"
	)

	updateRequest := UpdateTeam{
		Name:     &name,
		Type:     &teamType,
		Currency: &currency,
	}

	mux.HandleFunc("PUT /v1/teams/123", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodPut)
		v := new(UpdateTeam)
		err := json.NewDecoder(request.Body).Decode(v)
		require.NoError(t, err)

		assert.Equal(t, updateRequest, *v)

		writer.WriteHeader(http.StatusCreated)

		_, err = fmt.Fprint(writer, `{
					"id": 123,
					"name": "team"
			}
		`)
		require.NoError(t, err)
	})

	team, _, err := testClient.Teams.Update(t.Context(), 123, &updateRequest)
	require.NoError(t, err)

	assert.Equal(t, want, team)
}
