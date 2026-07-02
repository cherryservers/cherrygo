# cherrygo

Cherry Servers golang API client library for Cherry Servers RESTful API.

You can view the client API docs here: [https://pkg.go.dev/github.com/cherryservers/cherrygo/v3](https://pkg.go.dev/github.com/cherryservers/cherrygo/v3)

You can view Cherry Servers API docs here: [https://api.cherryservers.com/doc](https://api.cherryservers.com/doc)

## Table of Contents

- [cherrygo](#cherrygo)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Authentication](#authentication)
    - [Examples](#examples)
      - [Get teams](#get-teams)
      - [Get projects](#get-projects)
      - [Get plans](#get-plans)
      - [Get images](#get-images)
      - [Request new server](#request-new-server)
  - [License](#license)

## Installation

Add the library as a dependency to your project:
```
go get github.com/cherryservers/cherrygo/v3
```

### Authentication

To authenticate to the Cherry Servers API, you must have an API key. You can create API keys in the [Cherry Servers client portal](https://portal.cherryservers.com/settings/api-keys). API keys must be exported in the `CHERRY_API_KEY` environment variable or passed to the client directly.

Use an exported CHERRY_API_KEY environment variable:
```
export CHERRY_API_KEY="4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
```
```go
import (
	"context"
	"log"

	"github.com/cherryservers/cherrygo/v4"
)


func main() {
    ctx := context.Background()
    c, err := cherrygo.NewClient()
    if err != nil {
        log.Fatal(err)
    }
}
```
To pass a key to client without an environment variable:
```go
c, err := cherrygo.NewClient(cherrygo.WithAPIKey("your-api-key"))
```

### Examples

#### Get teams
You will need a team ID for subsequent function calls, for example, to get projects for a specified team, you will need to provide a team ID.
```go
teams, _, err := c.Teams.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, t := range teams {
    log.Printf("id: %d, name: %q, remaining promo credit: %f\n",
        t.ID, t.Name, t.Credit.Promo.Remaining)
}
```

#### Get projects
After you have your team ID, you can list your projects. You will need your project ID to list your servers or order new ones.
```go
projects, _, err := c.Projects.List(ctx, teamID, nil)
if err != nil {
    log.Fatal(err)
}

for _, p := range projects {
    log.Println(p.ID, p.Name)
}
```

#### Get plans
View available server plans.

```go
plans, _, err := c.Plans.List(ctx, teamID, nil)
if err != nil {
    log.Fatal(err)
}

for _, p := range plans {
    var hourlyPrice float32
    for _, pr := range p.Pricing {
        if pr.Unit == "Hourly" {
            hourlyPrice = pr.Price
        }
    }
    for _, r := range p.AvailableRegions {
        log.Printf("slug: %q, region: %q, stock: %d, type: %q, hourly price: %f",
            p.Slug, r.Slug, r.StockQty, p.Type, hourlyPrice)
    }
}
```

#### Get images
View OS images available for a specific plan.

```go
images, _, err := c.Images.List(ctx, planSlug, nil)
if err != nil {
    log.Fatal(err)
}

for _, i := range images {
    log.Println(i.Slug)
}
```

#### Request new server
```go
addServerRequest := cherrygo.CreateServer{
    ProjectID:   projectID,
    Image:       imageSlug,
    Region:      regionSlug,
    Plan:        planSlug,
}

server, _, err := c.Servers.Create(ctx, &addServerRequest)
if err != nil {
    log.Fatal(err)
}

log.Println(server.ID, server.Name, server.Hostname)
```

## License

See the [LICENSE](LICENSE.md) file for license rights and limitations.
