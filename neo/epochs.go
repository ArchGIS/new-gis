package neo

import "encoding/json"

type epochProps struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (db *DB) Epochs(params map[string]interface{}) ([]epochProps, error) {
	// stmt := fmt.Sprintf(epochsStatement, filterAuthors(params))
	buf, err := buildSingleStatementQuery(epochsStatement, params)
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

	return result.buildEpochsResponse()
}

const (
	epochsStatement = `
		MATCH (n:Epoch)
		RETURN n.id as id, n[{lang} + '_name'] as name
	`
)

func (resp *neoResponse) buildEpochsResponse() ([]epochProps, error) {
	data := resp.Results[0].Data
	var epochs []epochProps

	for _, row := range data {
		var epoch epochProps
		if err := json.Unmarshal(*row.Row, &epoch); err != nil {
			return nil, err
		}
		epochs = append(epochs, epoch)
	}

	return epochs, nil
}

func (row *epochProps) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&row.ID, &row.Name}
	return json.Unmarshal(buf, &tmp)
}
