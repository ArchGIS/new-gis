package neo

import (
	"strconv"
	"strings"

	"github.com/ArchGIS/new-gis/cypher"
	"github.com/gin-gonic/gin"
	"github.com/jmcvetta/neoism"
)

//////////////////////
////// Singular //////
//////////////////////

type (
	singleSite struct {
		Name     string   `json:"name"`
		Names    []string `json:"names"`
		Epoch    string   `json:"epoch"`
		Stype    string   `json:"type"`
		Cultures []string `json:"cultures"`

		ResCount  int     `json:"resCount"`
		ExcCount  int     `json:"excCount"`
		ExcArea   float64 `json:"excArea"`
		ArtiCount int     `json:"artiCount"`

		Heritages []nHeritage `json:"heritages"`

		LayersCount int     `json:"layersCount"`
		LayersTop   []layer `json:"layersTop"`
		LayersMid   []layer `json:"layersMid"`
		LayersBot   []layer `json:"layersBot"`

		Coords []siteSpatialReferences `json:"coords"`
	}

	layer struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name"`
		Epoch   string `json:"epoch"`
		Culture string `json:"culture"`
	}
)

const (
	coordPoint = iota + 1
	coordPolygon
)

func (db *DB) GetSite(id, lang string) (interface{}, error) {
	idInt, _ := strconv.Atoi(id)
	site := NewSite(uint64(idInt))

	var knowledges []siteNames
	var spatials []siteSpatialReferences
	var heritages []nHeritage
	var tEpoch []nEpoch
	var tSiteType []nSiteType
	var cultures []cultureNames
	var resCounts []researchCount
	var excCounts []excCount
	var artiCounts []artiCount

	cqs := []*neoism.CypherQuery{
		site.to(&knowledges),
		site.to(&spatials),
		site.to(&heritages),
		site.to(&tEpoch),
		site.to(&tSiteType),
		site.to(&cultures),
		site.to(&resCounts),
		site.to(&excCounts),
		site.to(&artiCounts),
	}

	err := db.CypherBatch(cqs)
	if err != nil {
		return nil, err
	}

	var response singleSite
	response.Names = knowledges[0].Names
	response.Coords = spatials
	response.Heritages = heritages
	response.Epoch = tEpoch[0].Name
	response.Stype = tSiteType[0].Name
	response.Cultures = cultures[0].Names
	response.ResCount = resCounts[0].Count
	response.ExcCount = excCounts[0].Count
	response.ExcArea = excCounts[0].Area
	response.ArtiCount = artiCounts[0].Count

	return response, nil
}

/*
 * Site researches
 */

func (db *DB) QuerySiteResearches(id, lang string) (interface{}, error) {
	return nil, nil
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

func (db *DB) Sites(req gin.H) (sites []pluralSite, err error) {
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

func siteFilterString(reqParams gin.H) string {
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
