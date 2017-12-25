package neo

import "encoding/json"

type (
	neoResponse struct {
		Results []struct {
			Data []struct {
				Row *json.RawMessage
			}
		}
	}

	neoQuery struct {
		Statements []statement `json:"statements"`
	}

	statement struct {
		Query  string                 `json:"statement"`
		Params map[string]interface{} `json:"parameters"`
	}

	point struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
)
