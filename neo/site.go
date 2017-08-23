package neo

import (
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/ArchGIS/new-gis/dbg"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

//////////////////////
////// Singular //////
//////////////////////

type singleSite struct {
	Names     []string `json:"names"`
	Epoch     int      `json:"epoch"`
	Stype     int      `json:"type"`
	Cultures  []string `json:"cultures"`
	ResCount  int      `json:"res_count"`
	ExcCount  int      `json:"exc_count"`
	ExcArea   int      `json:"exc_area"`
	ArtiCount int      `json:"arti_count"`
	HeritName string   `json:"herit_name"`
	HeritID   uint64   `json:"herit_id"`
}

const singleStatement = `
	MATCH (m:Monument {id: toInteger({id})})
	MATCH (m)-->(ep:Epoch)
	MATCH (m)-->(mt:MonumentType)
	MATCH (m)<--(k:Knowledge)
	MATCH (k)<--(r:Research)
	MATCH (k)-->(c:Culture)-[:translation {lang: {language}}]->(tr:Translate)
	OPTIONAL MATCH (r)-->(exc:Excavation)<--(m)
	OPTIONAL MATCH (exc)-->(a:Artifact)
	OPTIONAL MATCH (m)<--(h:Heritage)
	RETURN
		COLLECT(k.monument_name) AS names,
		ep.id AS epoch,
		mt.id AS type,
		COLLECT(tr.name) AS cultures,
		COUNT(r) AS res_count,
		COUNT(exc) AS exc_count,
		SUM(exc.area) AS exc_area,
		COUNT(a) AS arti_count,
		h.name AS herit_name,
		h.id AS herit_id
`

func (db *DB) GetSite(id, lang string) (*singleSite, error) {
	var dbResponse []singleSite
	cq := BuildCypherQuery(
		singleStatement,
		&dbResponse,
		neoism.Props{"id": id, "language": lang},
	)

	dbg.Dump(cq)
	err := db.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	dbg.Dump(dbResponse)

	return &dbResponse[0], nil
}

//////////////////////
////// Plural ////////
//////////////////////

type (
	item struct {
		Name     []string `json:"site_name"`
		ResName  []string `json:"research_name"`
		Epoch    int      `json:"epoch"`
		SiteType int      `json:"type"`
	}

	pluralSite struct {
		ID     uint64      `json:"id"`
		Item   item        `json:"item"`
		Coords coordinates `json:"coordinates"`
	}

	coordinates struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
)

const (
	pluralstatement = `
    MATCH (s:Monument)<--(k:Knowledge)
    MATCH (s)-[:has]->(st:MonumentType)
    MATCH (s)-[:has]->(e:Epoch)
    MATCH (r:Research)-[:has]->(k)
		%s
		WITH
			s.id as id,
			{
				site_name: collect(k.monument_name),
				research_name: collect(r.name),
				epoch: e.id,
				type: st.id
			} as item
    RETURN id, item
    SKIP {offset} LIMIT {limit}
	`

	entity = "Monument"
)

func (db *DB) Sites(req echo.Map) (sites []pluralSite, err error) {
	cq := BuildCypherQuery(
		cypher.Filter(pluralstatement, siteFilterString(req)),
		&sites,
		neoism.Props{
			"name":   cypher.BuildRegexpFilter(req["name"]),
			"epoch":  req["epoch"],
			"type":   req["type"],
			"offset": req["offset"],
			"limit":  req["limit"],
		},
	)

	err = db.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(sites) > 0 {
		ids := make([]uint64, len(sites))
		coords := make([]coordinates, len(sites))

		for i, v := range sites {
			ids[i] = v.ID
		}

		coordsCQ := BuildCypherQuery(
			cypher.BuildCoordinates(ids, entity, false),
			&coords,
			neoism.Props{},
		)

		err = db.Cypher(&coordsCQ)
		if err != nil {
			return nil, err
		}

		for i := range coords {
			sites[i].Coords = coords[i]
		}
	}

	return sites, nil
}

func siteFilterString(reqParams echo.Map) string {
	var filter []string
	var stmt string

	if reqParams["name"] != "" {
		filter = append(filter, "k.monument_name =~ {name}")
	}
	if reqParams["epoch"] != 0 {
		filter = append(filter, "e.id = {epoch}")
	}
	if reqParams["type"] != 0 {
		filter = append(filter, "st.id = {type}")
	}
	if len(filter) > 0 {
		stmt = "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}
