package neo

import (
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/neoism"
)

type (
	item struct {
		ResName    string `json:"research_name"`
		ReportName string `json:"report_name"`
		Year       int64  `json:"year"`
		Author     string `json:"author_name"`
		Type       string `json:"res_type"`
	}

	pluralResearch struct {
		ID     uint64        `json:"id"`
		Item   item          `json:"item"`
		Coords []coordinates `json:"coordinates"`
	}

	resCoord struct {
		Rows []coordinates `json:"rows"`
	}
	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
)

const (
	pluralResearchStatement = `
		MATCH (n:Research)-[:hasauthor]->(a:Author)
		MATCH (n)-[:has]->(rep:Report)
		MATCH (n)-[:has]->(rtype:ResearchType)-[:translation {lang: {language}}]->(trans:Translate)
		%s
		WITH n.id as id, {
			research_name: n.name,
			report_name: rep.name,
			year: n.year,
			author_name: a.name,
			res_type: trans.name
		} as item
		RETURN id, item
		SKIP {offset} LIMIT {limit}
	`

	entity = "Research"
)

func (db *DB) Researches(req gin.H) (res []pluralResearch, err error) {
	cq := BuildCypherQuery(
		cypher.Filter(statement, researchFilterString(req)),
		&res,
		neoism.Props{
			"language": req["lang"],
			"name":     cypher.BuildRegexpFilter(req["name"]),
			"year":     req["year"],
			"offset":   req["offset"],
			"limit":    req["limit"],
		},
	)

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		ids := make([]uint64, len(res))
		for i, v := range res {
			ids[i] = v.ID
		}

		listResearchSites := []cypher.ListIdsOfSites{}
		cq = BuildCypherQuery(
			cypher.IdsOfResearchSites(ids),
			&listResearchSites,
			neoism.Props{},
		)
		err = db.Cypher(&cq)
		if err != nil {
			return nil, err
		}

		coords := make([]resCoord, len(res))

		coordsCQ := BuildCypherQuery(
			cypher.ResearchCoords(listResearchSites),
			&coords,
			neoism.Props{},
		)

		err = db.Cypher(&coordsCQ)
		if err != nil {
			return nil, err
		}

		for i := range coords {
			res[i].Coords = coords[i].Rows
		}
	}

	return res, nil
}

func researchFilterString(reqParams gin.H) string {
	var filter []string
	var stmt string

	if reqParams["name"] != "" {
		filter = append(filter, "n.name =~ {name}")
	}
	if reqParams["year"] != MinInt {
		filter = append(filter, "n.year = {year}")
	}
	if len(filter) > 0 {
		stmt = "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}
