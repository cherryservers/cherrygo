# README #

Cherry Servers golang API for Cherry Servers RESTful API.

Installation
------------

### Examples ###


```go
	c, err := cherrygo.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	teams, _, err := c.Teams.List()
	if err != nil {
		log.Fatal("Error", err)
	}

	for _, t := range teams {

		fmt.Fprintf("%v\t%v\t%v\t%v\t%v\n",
			t.ID, t.Name, t.Credit.Promo.Remaining, t.Credit.Promo.Usage, t.Credit.Resources.Pricing.Price)
	}
```