package couch

import (
	"encoding/json"
	"fmt"

	"github.com/byuoitav/av-control-api/api/base"
)

type buildingQueryResponse struct {
	Docs     []building `json:"docs"`
	Bookmark string     `json:"bookmark"`
	Warning  string     `json:"warning"`
}

type building struct {
	Rev string `json:"_rev,omitempty"`
	*base.Building
}

// GetAllBuildings returns all the buildings found in the database
func (c *CouchDB) GetAllBuildings() ([]base.Building, error) {
	var toReturn []base.Building
	var query IDPrefixQuery
	query.Selector.ID.GT = "\x00"
	query.Limit = 1000

	b, err := json.Marshal(query)
	if err != nil {
		return toReturn, fmt.Errorf("failed to marshal query to get all buildings: %s", err)
	}

	var resp buildingQueryResponse

	err = c.MakeRequest("POST", fmt.Sprintf("%v/_find", BUILDINGS), "application/json", b, &resp)
	if err != nil {
		return toReturn, fmt.Errorf("failed to get all buildings: %s", err)
	}

	for _, doc := range resp.Docs {
		toReturn = append(toReturn, *doc.Building)
	}

	return toReturn, err
}
