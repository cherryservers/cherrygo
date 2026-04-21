package cherrygo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImages_List(t *testing.T) {
	setup()
	defer teardown()

	expected := []Image{
		{
			ID:   1,
			Name: "CloudLinux 7 64bit",
			Slug: "cloudlinux_7",
			Pricing: []Pricing{
				{
					Price:    0.015,
					Taxed:    false,
					Currency: "EUR",
					Unit:     "Hourly",
				},
			},
		},
		{
			ID:   2,
			Name: "Ubuntu 20.04 64bit",
			Slug: "ubuntu_20_04",
		},
	}

	mux.HandleFunc("/v1/plans/e5_1620v4/images", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		_, err := fmt.Fprint(writer, `[
			{
				"id": 1,
				"name": "CloudLinux 7 64bit",
				"slug": "cloudlinux_7",
				"pricing": [
					{
						"price": 0.015,
						"taxed": false,
						"currency": "EUR",
						"unit": "Hourly"
					}
				]
			},
			{
				"id": 2,
				"name": "Ubuntu 20.04 64bit",
				"slug": "ubuntu_20_04"
			}
		]`)

		require.NoError(t, err)
	})

	images, _, err := testClient.Images.List(t.Context(), "e5_1620v4", nil)
	if err != nil {
		t.Errorf("Images.List returned %+v", err)
	}

	if !reflect.DeepEqual(images, expected) {
		t.Errorf("Images.List returned %+v, expected %+v", images, expected)
	}
}
