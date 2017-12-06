package neo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type (
	authorsProps struct {
		ID       uint64   `json:"id"`
		Name     string   `json:"author_name"`
		ResNames []string `json:"research_names"`
	}

	neoQuery struct {
		Statements []Statement `json:"statements"`
	}

	Statement struct {
		Query  string                 `json:"statement"`
		Params map[string]interface{} `json:"parameters"`
	}
)

func (db *DB) Authors(params map[string]interface{}) (interface{}, error) {
	stmt := fmt.Sprintf(
		"MATCH (a:Author)<-[:hasauthor]-(r:Research) "+
			"%s "+
			"WITH a.id as id, a.name as author_name, collect(r.name) as research_names "+
			"RETURN id, author_name, research_names "+
			"SKIP {offset} LIMIT {limit} ",
		filterAuthors(params),
	)
	reqB, err := json.Marshal(
		neoQuery{
			Statements: []Statement{
				Statement{Query: stmt, Params: params},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(reqB)
	req, err := http.NewRequest("POST", os.Getenv("NEO4J_REST_TR_COMMIT"), buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(os.Getenv("NEO4J_USER"), os.Getenv("NEO4J_PASSWORD"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	// for i := range result.results

	return result.(map[string]interface{})["results"], nil
}

func filterAuthors(req map[string]interface{}) (filter string) {
	if req["name"] != "" {
		filter = "WHERE a.name CONTAINS {name}"
	}

	return filter
}
