package neo

import "github.com/jmcvetta/neoism"
import "fmt"

type (
	site struct {
		ID uint64
	}

	nodeProps map[string]interface{}
)

func NewSite(id uint64) *site {
	return &site{ID: id}
}

func (siteObj *site) to(label string, result interface{}, props []string) *neoism.CypherQuery {
	cq := neoism.CypherQuery{
		Statement: fmt.Sprintf(siteToKnowledge, siteObj.ID),
		Result:    result,
	}

	return &cq
}

type knowledge struct {
	ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToKnowledge = `
	MATCH (s:Monument {id: %d})<--(k:Knowledge)
	RETURN
		k.id as id,
		k {.monument_name} as data
`

type spatialReference struct {
	ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToSpatial = `
	MATCH (s:Monument {id: %d})-->(n:SpatialReference)
	RETURN
		n.id as ID,
		n {.x, .y} as data
`

type nSiteType struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
}

const siteToType = `
	MATCH (s:Monument {id: %d})-->(n:MonumentType)-[:translation {lang: "%s"}]->(tr:Translate)
	RETURN
		n.id as id,
		tr.name as name
`

type nEpoch struct {
	ID   uint16 `json:"id"`
	Name string `json:"name"`
}

const siteToEpoch = `
MATCH (s:Monument {id: %d})-->(n:Epoch)-[:translation {lang: "%s"}]->(tr:Translate)
RETURN
	n.id as id,
	tr.name as name
`

type Excavation struct {
	ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToExc = `
	MATCH (s:Monument {id: %d})-->(n:Excavation)
	RETURN
		n.id as ID,
		n {.name} as data
`
