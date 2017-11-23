package neo

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func (db *DB) Graphql(req []byte) (interface{}, error) {
	result, err := db.graphql(req)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DB) graphql(query []byte) (interface{}, error) {
	buf := bytes.NewBuffer(query)
	req, err := http.NewRequest("POST", os.Getenv("NEO4J_GRAPHQL_URL"), buf)
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

	return result.(map[string]interface{})["data"], nil
}
