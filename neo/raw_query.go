package neo

import "encoding/json"

func (db *DB) RawQuery(query string) (interface{}, error) {
	buf, err := buildSingleStatementQuery(query, map[string]interface{}{})
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

	return extractRawRows(&result)
}

func extractRawRows(resp *neoResponse) ([]interface{}, error) {
	data := resp.Results[0].Data
	var results []interface{}
	
	for _, row := range data {
		var rowData interface{}
		if err := json.Unmarshal(*row.Row, &rowData); err != nil {
			return nil, err
		}
		results = append(results, rowData)
	}

	return results, nil
}