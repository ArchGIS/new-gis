package neo

import (
	"fmt"
	"strings"

	"github.com/jmcvetta/neoism"
)

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

func (siteObj *site) to(result interface{}) *neoism.CypherQuery {
	var query string

	switch result.(type) {
	case *[]siteNames:
		query = fmt.Sprintf(siteToKnowledge, siteObj.ID)
	case *[]siteSpatialReferences:
		query = fmt.Sprintf(siteToSpatial, siteObj.ID)
	case *[]nHeritage:
		query = fmt.Sprintf(siteToHeritage, siteObj.ID)
	case *[]nEpoch:
		query = fmt.Sprintf(siteToEpoch, siteObj.ID, "en")
	case *[]nSiteType:
		query = fmt.Sprintf(siteToType, siteObj.ID, "en")
	case *[]cultureNames:
		query = fmt.Sprintf(siteToCultures, siteObj.ID, "en")
	case *[]researchCount:
		query = fmt.Sprintf(siteToResCount, siteObj.ID)
	case *[]excCount:
		query = fmt.Sprintf(siteToExcCount, siteObj.ID)
	case *[]artiCount:
		query = fmt.Sprintf(siteToArtiCount, siteObj.ID)
	}
	return &neoism.CypherQuery{
		Statement: query,
		Result:    result,
	}
}

type siteNames struct {
	Names []string `json:"names"`
}

const siteToKnowledge = `
	MATCH (:Monument {id: %d})<--(k:Knowledge)
	RETURN COLLECT(k.monument_name) as names
`

type siteSpatialReferences struct {
	Date     uint64  `json:"date"`
	Accuracy int     `json:"accuracy"`
	Points   []point `json:"points"`
}

const siteToSpatial = `
	MATCH (:Monument {id: %d})-->(sp:SpatialReference)-->(spt:SpatialReferenceType)
	WITH sp, spt
	ORDER BY spt.id ASC, sp.date DESC
	RETURN
		[{x: sp.x, y: sp.y}] as points,
		spt.id as accuracy,
		sp.date as date
`

type nSiteType struct {
	Name string `json:"name"`
}

const siteToType = `
	MATCH (:Monument {id: %d})-->(:MonumentType)-[:translation {lang: "%s"}]->(tr:Translate)
	RETURN tr.name as name
`

type nEpoch struct {
	Name string `json:"name"`
}

const siteToEpoch = `
	MATCH (:Monument {id: %d})-->(:Epoch)-[:translation {lang: "%s"}]->(tr:Translate)
	RETURN tr.name as name
`

type excCount struct {
	Count int     `json:"count"`
	Area  float64 `json:"area"`
}

const siteToExcCount = `
	MATCH (:Monument {id: %d})-->(n:Excavation)
	RETURN
		COUNT(n) as count,
		SUM(n.area) as area
`

type nHeritage struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

const siteToHeritage = `
	MATCH (:Monument {id: %d})<--(n:Heritage)
	RETURN
		n.id as id,
		n.name as name
`

type cultureNames struct {
	Names []string `json:"names"`
}

const siteToCultures = `
	MATCH (:Monument {id: %d})<--(:Knowledge)-->(:Culture)-[:translation {lang: "%s"}]->(tr:Translate)
	RETURN COLLECT(tr.name) as names
`

type researchCount struct {
	Count int `json:"count"`
}

const siteToResCount = `
	MATCH (:Monument {id: %d})<--(:Knowledge)<--(r:Research)
	RETURN COUNT(r) as count
`

type artiCount struct {
	Count int `json:"count"`
}

const siteToArtiCount = `
	MATCH (:Monument {id: %d})-->(:Excavation)-->(a:Artifact)
	RETURN COUNT(a) as count
`
