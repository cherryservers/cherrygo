package cherrygo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUser_Current(t *testing.T) {
	setup()
	defer teardown()

	expected := User{
		ID:                    123,
		FirstName:             "Ei",
		LastName:              "Jei",
		Email:                 "email@email.com",
		EmailVerified:         true,
		Phone:                 "37060000000",
		SecurityPhone:         "37060000000",
		SecurityPhoneVerified: false,
	}

	mux.HandleFunc("/v1/user", func(writer http.ResponseWriter, request *http.Request) {
		testMethod(t, request, http.MethodGet)
		fmt.Fprint(writer, `{
				"id":123,
				"first_name": "Ei",
				"last_name": "Jei",
				"email": "email@email.com",
				"email_verified": true,
				"phone": "37060000000",
				"security_phone": "37060000000",
				"security_phone_verified": false
			 }`)
	})

	user, _, err := client.Users.CurrentUser(nil)
	if err != nil {
		t.Errorf("Users.CurrentUser returned %+v", err)
	}

	if !reflect.DeepEqual(user, expected) {
		t.Errorf("Users.CurrentUser returned %+v, expected %+v", user, expected)
	}
}
