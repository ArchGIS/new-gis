package neo

import (
	"encoding/json"
	"fmt"
)

type (
	authorsProps struct {
		ID       uint64   `json:"id"`
		Name     string   `json:"author_name"`
		ResNames []string `json:"researches_names"`
	}
)

const authorsCypherQuery = "MATCH (a:Author)<-[:hasauthor]-(r:Research) " +
	"%s " +
	"WITH a.id as id, a.name as author_name, collect(r.name) as research_names " +
	"RETURN id, author_name, research_names " +
	"SKIP {offset} LIMIT {limit} "

// Authors returning data from database about authors
func (db *DB) Authors(params map[string]interface{}) (interface{}, error) {
	stmt := fmt.Sprintf(authorsCypherQuery, filterAuthors(params))
	buf, err := buildSingleStatementQuery(stmt, params)
	if err != nil {
		return nil, err
	}

	req, err := buildNeoRequest(buf)
	if err != nil {
		return nil, err
	}

	resp, err := neoClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result neoResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.buildAuthorsResponse()
}

func filterAuthors(req map[string]interface{}) (filter string) {
	if req["name"] != "" {
		filter = "WHERE a.name CONTAINS {name}"
	}

	return filter
}

func (resp *neoResponse) buildAuthorsResponse() ([]authorsProps, error) {
	data := resp.Results[0].Data
	var authors []authorsProps

	for _, row := range data {
		var author authorsProps
		if err := json.Unmarshal(*row.Row, &author); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}

	return authors, nil
}

func (row *authorsProps) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&row.ID, &row.Name, &row.ResNames}
	return json.Unmarshal(buf, &tmp)
}
