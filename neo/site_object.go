package neo

import "github.com/jmcvetta/neoism"
import "fmt"
import "strings"

type (
	site struct {
		ID uint64
	}

	nodeProps map[string]interface{}
)

func NewSite(id uint64) *site {
	return &site{ID: id}
}

func toString(props []string) string {
	projected := make([]string, len(props))
	for i, v := range props {
		projected[i] = "." + v
	}
	return strings.Join(projected, ",")
}

func (siteObj *site) to(result interface{}, props []string) *neoism.CypherQuery {
	var query string
	var projection string

	switch result.(type) {
	case *[]knowledge:
		query = siteToKnowledge
		projection = toString(props)
	}
	return &neoism.CypherQuery{
		Statement: fmt.Sprintf(query, siteObj.ID, projection),
		Result:    result,
	}
}

type knowledge struct {
	ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToKnowledge = `
	MATCH (s:Monument {id: %d})<--(k:Knowledge)
	RETURN
		k.id as id,
		k {%s} as data
`

type spatialReference struct {
	// ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToSpatial = `
	MATCH (s:Monument {id: %d})-->(sp:SpatialReference)-->(spt:SpatialReferenceType)
	WITH sp, spt
	ORDER BY srt.id ASC, sr.date DESC
	RETURN {
		x: sp.x,
		y: sp.y,
		accuracy: spt.id
	} as data
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

type excavation struct {
	ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToExc = `
	MATCH (s:Monument {id: %d})-->(n:Excavation)
	RETURN
		n.id as ID,
		n {.name} as data
`

type nHeritage struct {
	ID   uint64    `json:"id"`
	Data nodeProps `json:"data"`
}

const siteToHeritage = `
	MATCH (s:Monument {id: %d})<--(n:Heritage)
	RETURN
		n.id as ID,
		n {.name} as data
`
