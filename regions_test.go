package cherrygo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRegions_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []Region{
		{
			ID:         1,
			Name:       "EU-Nord-1",
			Slug:       "eu_nord_1",
			RegionIso2: "LT",
			Href:       "/regions/1",
		},
		{
			ID:         2,
			Name:       "EU-West-1",
			Slug:       "eu_west_1",
			RegionIso2: "NL",
			Href:       "/regions/2",
		},
	}

	mux.HandleFunc("/v1/regions", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		fmt.Fprint(writer, `[
			{
				"id":1,
				"name":"EU-Nord-1",
				"slug":"eu_nord_1",
				"region_iso_2":"LT",
				"href":"/regions/1"
			 },
			 {
				"id":2,
				"name":"EU-West-1",
				"slug":"eu_west_1",
				"region_iso_2":"NL",
				"href":"/regions/2"
			 }
		]`)
	})

	regions, _, err := client.Regions.List(nil)
	if err != nil {
		t.Errorf("Regions.List returned %+v", err)
	}

	if !reflect.DeepEqual(regions, expected) {
		t.Errorf("Regions.List returned %+v, expected %+v", regions, expected)
	}
}

func TestRegion_Get(t *testing.T) {
	setup()
	defer teardown()

	expected := Region{
		ID:         1,
		Name:       "EU-Nord-1",
		Slug:       "eu_nord_1",
		RegionIso2: "LT",
		Href:       "/regions/1",
	}

	mux.HandleFunc("/v1/regions/eu_nord_1", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		fmt.Fprint(writer, `{
			"id":1,
			"name":"EU-Nord-1",
			"slug":"eu_nord_1",
			"region_iso_2":"LT",
			"href":"/regions/1"
		}`)
	})

	region, _, err := client.Regions.Get("eu_nord_1", nil)
	if err != nil {
		t.Errorf("Regions.Get returned %+v", err)
	}

	if !reflect.DeepEqual(region, expected) {
		t.Errorf("Regions.Get returned %+v, expected %+v", region, expected)
	}
}
