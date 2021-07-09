package cherrygo

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestPlans_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/teams/"+strconv.Itoa(teamID)+"/plans", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		response := `[{"id":625,"name":"Cloud VPS 1","title":"Cloud VPS 1","custom":false,"category":"Shared resources","softwares":[{"image":{"name":"Ubuntu 18.04 64bit"}}],"specs":{"cpus":{"count":1,"name":"Cloud VPS 1","cores":1,"frequency":0.0,"unit":"GHz"},"memory":{"count":1,"total":1,"unit":"GB","name":"1GB"},"storage":[{"count":1,"name":"20GB SSD","size":20,"unit":"GB","type":"SSD"}],"nics":{"name":"1Gbps"},"bandwidth":{"name":"1TB"}},"pricing":[{"price":0.015,"currency":"EUR","taxed":false,"unit":"Hourly","id":37}],"available_regions":[{"id":1,"name":"EU-Nord-1","region_iso_2":"LT","stock_qty":122,"spot_qty":0,"location":"Lithuania, Vilnius"},{"id":2,"name":"EU-West-1","region_iso_2":"NL","stock_qty":91,"spot_qty":0,"location":"Netherlands, Amsterdam"}]}]`
		fmt.Fprint(writer, response)
	})

	plans, _, err := client.Plans.List(teamID, nil)
	if err != nil {
		t.Errorf("Plans.List returned %+v", err)
	}

	expected := []Plans{
		{
			ID:     625,
			Name:   "Cloud VPS 1",
			Custom: false,
			Specs: Specs{
				Cpus: Cpus{
					Count:     1,
					Name:      "Cloud VPS 1",
					Cores:     1,
					Frequency: 0.0,
					Unit:      "GHz",
				},
				Memory: Memory{
					Count: 1,
					Total: 1,
					Unit:  "GB",
					Name:  "1GB",
				},
				Storage: []Storage{{
					Count: 1,
					Name:  "20GB SSD",
					Size:  20,
					Unit:  "GB",
				}},
				//Raid: Raid{},
				Nics: Nics{
					Name: "1Gbps",
				},
				Bandwidth: Bandwidth{
					Name: "1TB",
				},
			},
			Pricing: []Pricing{{
				Price:    0.015,
				Taxed:    false,
				Currency: "EUR",
				Unit:     "Hourly",
			}},
			AvailableRegions: []AvailableRegions{
				{
					ID:         1,
					Name:       "EU-Nord-1",
					RegionIso2: "LT",
					StockQty:   122,
				},
				{
					ID:         2,
					Name:       "EU-West-1",
					RegionIso2: "NL",
					StockQty:   91,
				},
			},
		},
	}

	if !reflect.DeepEqual(plans, expected) {
		t.Errorf("Plans.List  plans returned %+v, expected %+v", plans, expected)
	}
}
