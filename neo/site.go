package neo

import (
	"strconv"
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/jmcvetta/neoism"
	"github.com/labstack/echo"
)

//////////////////////
////// Singular //////
//////////////////////

type (
	singleSite struct {
		Name     string   `json:"name"`
		Names    []string `json:"names"`
		Epoch    int      `json:"epoch"`
		Stype    int      `json:"type"`
		Cultures []string `json:"cultures"`

		ResCount  int `json:"resCount"`
		ExcCount  int `json:"excCount"`
		ExcArea   int `json:"excArea"`
		ArtiCount int `json:"artiCount"`

		Heritages []heritage `json:"heritages"`

		LayersCount int     `json:"layersCount"`
		LayersTop   []layer `json:"layersTop"`
		LayersMid   []layer `json:"layersMid"`
		LayersBot   []layer `json:"layersBot"`

		Coords []siteCoords `json:"coords"`
	}

	heritage struct {
		Name string `json:"name"`
		ID   uint64 `json:"id"`
	}

	layer struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name"`
		Epoch   string `json:"epoch"`
		Culture string `json:"culture"`
	}

	siteCoords struct {
		Date     uint64  `json:"date"`
		Type     int     `json:"type"`
		Accuracy int     `json:"accuracy"`
		Points   []point `json:"points"`
		// Actual   bool    `json:"actual"`
	}
)

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
	OPTIONAL MATCH (m)-[:has]->(sr:SpatialReference)-[:has]->(srt:SpatialReferenceType)
	WITH sr, srt, m, ep, mt, k, r, c, tr, exc, a, h
	ORDER BY srt.id ASC, sr.date DESC
	RETURN
		COLLECT(k.monument_name) AS names,
		ep.id AS epoch,
		mt.id AS type,
		COLLECT(tr.name) AS cultures,
		COUNT(r) AS resCount,
		COUNT(exc) AS excCount,
		SUM(exc.area) AS excArea,
		COUNT(a) AS artiCount,
		COLLECT({name: h.name, id: h.id}) as heritages,
		COLLECT({date: sr.date, accuracy: srt.id, type: 1, points: [{x: sr.x, y: sr.y}]}) AS coords
`

func (db *DB) GetSite(id, lang string) ([]knowledge, error) {
	// var dbResponse []singleSite
	// cq := BuildCypherQuery(
	// 	singleStatement,
	// 	&dbResponse,
	// 	neoism.Props{"id": id, "language": lang},
	// )

	// err := db.Cypher(&cq)
	// if err != nil {
	// 	return nil, err
	// }
	// dbg.Dump(dbResponse)

	idInt, _ := strconv.Atoi(id)
	site := NewSite(uint64(idInt))

	var knowledges []knowledge
	cqs := []*neoism.CypherQuery{
		site.to("", &knowledges, []string{}),
	}

	err := db.CypherBatch(cqs)
	if err != nil {
		return nil, err
	}
	return knowledges, nil //&dbResponse[0], nil
}

//////////////////////
////// Plural ////////
//////////////////////

type (
	siteItem struct {
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

	monument = "Monument"
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
			cypher.BuildCoordinates(ids, monument, false),
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
