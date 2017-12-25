package neo

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

var neoClient *http.Client

func init() {
	neoClient = &http.Client{}
}

func buildSingleStatementQuery(stmt string, params map[string]interface{}) (io.Reader, error) {
	query, err := json.Marshal(
		neoQuery{
			Statements: []Statement{
				Statement{Query: stmt, Params: params},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(query), nil
}

func buildNeoRequest(buf io.Reader) (*http.Request, error) {
	req, err := http.NewRequest("POST", os.Getenv("NEO4J_REST_TR_COMMIT"), buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(os.Getenv("NEO4J_USER"), os.Getenv("NEO4J_PASSWORD"))

	return req, nil
}
