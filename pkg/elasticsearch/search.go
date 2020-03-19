package elasticsearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/vshn/waf-tool/pkg/model"
)

// SearchUniqueID to search for a ModSecurity unique ID
func (c *client) SearchUniqueID(uniqueID string) (model.SearchResult, error) {

	var buf bytes.Buffer
	type m map[string]interface{}
	query := m{
		"query": m{
			"bool": m{
				"should": []m{{
					"match": m{
						"apache-access.uniqueID": uniqueID,
					}}, {
					"match": m{
						"modsec-alert.uniqueID": uniqueID,
					}},
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return model.SearchResult{}, fmt.Errorf("error encoding query %w", err)
	}

	res, err := c.es.Search(
		c.es.Search.WithHeader(map[string]string{"Authorization": "Bearer " + c.token}),
		c.es.Search.WithIndex(indexWildcard),
		c.es.Search.WithIgnoreUnavailable(true),
		c.es.Search.WithBody(&buf),
	)
	defer res.Body.Close()
	if err != nil {
		return model.SearchResult{}, err
	}
	if res.IsError() {
		var e model.ErrorResponse
		if res.StatusCode == http.StatusUnauthorized {
			return model.SearchResult{}, errors.New("401 unauthorized")
		}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return model.SearchResult{}, fmt.Errorf("failure at parsing the response body: %w", err)
		}
		return model.SearchResult{}, fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e.Error.Type,
			e.Error.Reason,
		)
	}
	var searchResult model.SearchResult
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return model.SearchResult{}, err
	}
	return searchResult, nil
}
