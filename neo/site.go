package neo

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
	// idInt, _ := strconv.Atoi(id)
	// site := NewSite(uint64(idInt))

	// var knowledges []siteNames
	// var spatials []siteSpatialReferences
	// var heritages []nHeritage
	// var tEpoch []nEpoch
	// var tSiteType []nSiteType
	// var cultures []cultureNames
	// var resCounts []researchCount
	// var excCounts []excCount
	// var artiCounts []artiCount

	// cqs := []*neoism.CypherQuery{
	// 	site.to(&knowledges),
	// 	site.to(&spatials),
	// 	site.to(&heritages),
	// 	site.to(&tEpoch),
	// 	site.to(&tSiteType),
	// 	site.to(&cultures),
	// 	site.to(&resCounts),
	// 	site.to(&excCounts),
	// 	site.to(&artiCounts),
	// }

	// err := db.CypherBatch(cqs)
	// if err != nil {
	// 	return nil, err
	// }

	// var response singleSite
	// response.Names = knowledges[0].Names
	// response.Coords = spatials
	// response.Heritages = heritages
	// response.Epoch = tEpoch[0].Name
	// response.Stype = tSiteType[0].Name
	// response.Cultures = cultures[0].Names
	// response.ResCount = resCounts[0].Count
	// response.ExcCount = excCounts[0].Count
	// response.ExcArea = excCounts[0].Area
	// response.ArtiCount = artiCounts[0].Count

	return nil, nil
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
	pluralSite struct {
		ID     int64    `json:"id"`
		Item   siteItem `json:"item"`
		Coords coordsT  `json:"coordinates"`
	}

	siteItem struct {
		Names    interface{} `json:"site_names"`
		ResNames interface{} `json:"research_names"`
		Epoch    int64       `json:"epoch"`
		SiteType int64       `json:"type"`
	}

	coordsT struct {
		X interface{} `json:"x"`
		Y interface{} `json:"y"`
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
			collect(k.monument_name) as site_name,
			collect(r.name) as research_name,
			e.id as epoch,
			st.id as type
    RETURN id, site_name, research_name, epoch, type
    SKIP {offset} LIMIT {limit}
	`

	monument = "Monument"
)

func idsForQueriengCoordinates(ids []int64) gin.H {
	params := make(gin.H)
	for i, v := range ids {
		params["id"+strconv.Itoa(i)] = v
	}

	return params
}

func (db *DB) Sites(req gin.H) ([]pluralSite, error) {
	rows, err := db.QueryNeo(
		fmt.Sprintf(pluralstatement, siteFilterString(req)),
		gin.H{
			"name":   buildRegexpFilter(req["name"]),
			"epoch":  req["epoch"],
			"type":   req["type"],
			"offset": req["offset"],
			"limit":  req["limit"],
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not query in neo4j: %v", err)
	}

	var sites []pluralSite
	var ids []int64
	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("didn't get rows: %v", err)
		} else if err == io.EOF {
			continue
		}

		sites = append(sites, pluralSite{
			ID: row[0].(int64),
			Item: siteItem{
				Names:    row[1],
				ResNames: row[2],
				Epoch:    row[3].(int64),
				SiteType: row[4].(int64),
			},
		})
		ids = append(ids, row[0].(int64))
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("error when closing statement: %v", err)
	}

	preStmt := actualCoordinates(len(ids), monument)
	coordRows, err := db.QueryNeo(preStmt, idsForQueriengCoordinates(ids))
	if err != nil {
		return nil, fmt.Errorf("could not query in neo4j: %v", err)
	}
	defer func() {
		err = coordRows.Close()
		if err != nil {
			fmt.Printf("closing coordinates statement failed: %v", err)
		}
	}()

	coords, _, err := coordRows.All()
	if err != nil {
		return nil, fmt.Errorf("could not get rows of data: %v", err)
	}

	for i := range sites {
		sites[i].Coords.X = coords[i][0]
		sites[i].Coords.Y = coords[i][1]
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
