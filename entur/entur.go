package entur

import (
	"context"

	"github.com/shurcooL/graphql"
)

// Internal types

type stopPlace struct {
	ID   graphql.String
	Name struct {
		Value graphql.String
	}
}

// External types

// Stop which can receive arrivals
type Stop struct {
	Name string
	ID   string
}

// GetStop Fetches a stop based on its name
func GetStop(name string) (Stop, error) {
	var stop Stop

	client := graphql.NewClient("https://api.entur.org/stop_places/1.0/graphql", nil)

	var query struct {
		StopPlace []stopPlace `graphql:"stopPlace(query: $stop)"`
	}
	variables := map[string]interface{}{
		"stop": graphql.String(name),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return stop, err
	}

	return Stop{
		string(query.StopPlace[0].Name.Value),
		string(query.StopPlace[0].ID),
	}, nil
}
