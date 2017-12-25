package neo

import "encoding/json"

type (
	siteTypeProps struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

func (db *DB) SiteTypes(params map[string]interface{}) ([]siteTypeProps, error) {
	buf, err := buildSingleStatementQuery(siteTypesStatement, params)
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

	return result.siteTypes()
}

const (
	siteTypesStatement = `
		MATCH (n:MonumentType)
		RETURN n.id as id, n[{lang} + '_name'] as name
	`
)

func (resp *neoResponse) siteTypes() ([]siteTypeProps, error) {
	data := resp.Results[0].Data
	var types []siteTypeProps

	for _, row := range data {
		var site siteTypeProps
		if err := json.Unmarshal(*row.Row, &site); err != nil {
			return nil, err
		}
		types = append(types, site)
	}

	return types, nil
}

func (row *siteTypeProps) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&row.ID, &row.Name}
	return json.Unmarshal(buf, &tmp)
}