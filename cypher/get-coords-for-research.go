package cypher

import (
	"fmt"
	"strings"
)

type (
	ListIdsOfSites struct {
		IDs []uint64 `json:"ids"`
	}
)

func IdsOfResearchSites(ids []uint64) string {
	stmt := `MATCH (r:Research {id: %d})-->(k:Knowledge)-->(m:Monument)
		WITH r, COLLECT(m) AS sites
		RETURN EXTRACT(x IN sites | x.id) AS ids`

	acc := make([]string, len(ids))
	for i, v := range ids {
		acc[i] = fmt.Sprintf(stmt, v)
	}

	return strings.Join(acc, " UNION ALL ")
}

func ResearchCoords(listIds []ListIdsOfSites) string {
	acc := make([]string, len(listIds))

	for i, v := range listIds {
		acc[i] = BuildCoordinates(v.IDs, "Monument", true)
	}

	return strings.Join(acc, " UNION ALL ")
}
