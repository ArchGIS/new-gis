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
		ResNames []string `json:"researches_names"`
	}

	neoResponse struct {
		Results []struct {
			Data []struct {
				Row *json.RawMessage
			}
		}
	}

	neoQuery struct {
		Statements []Statement `json:"statements"`
	}

	Statement struct {
		Query  string                 `json:"statement"`
		Params map[string]interface{} `json:"parameters"`
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
