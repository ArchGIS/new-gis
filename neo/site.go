package neo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver/encoding"
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

		ResCount  int64   `json:"resCount"`
		ExcCount  int64   `json:"excCount"`
		ExcArea   float64 `json:"excArea"`
		ArtiCount int64   `json:"artiCount"`

		Heritages []*nHeritage `json:"heritages"`

		LayersCount int64   `json:"layersCount"`
		LayersTop   []layer `json:"layersTop"`
		LayersMid   []layer `json:"layersMid"`
		LayersBot   []layer `json:"layersBot"`

		Coords []*siteSpatialReferences `json:"coords"`
	}

	layer struct {
		ID      uint64 `json:"id"`
		Name    string `json:"name"`
		Epoch   string `json:"epoch"`
		Culture string `json:"culture"`
	}

	siteSpatialReferences struct {
		Date     int64   `json:"date"`
		Accuracy int64   `json:"accuracy"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}

	nHeritage struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}
)

const (
	coordPoint = iota + 1
	coordPolygon
)

func (db *DB) GetSite(req map[string]interface{}) (*singleSite, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	var response singleSite
	names, err := db.getSiteNames(params)
	if err != nil {
		return nil, fmt.Errorf("names query failed: %v", err)
	}
	response.Names = names

	coords, err := db.getSiteCoordinates(params)
	if err != nil {
		return nil, fmt.Errorf("coords query failed: %v", err)
	}
	response.Coords = coords

	siteType, err := db.getSiteType(params)
	if err != nil {
		return nil, fmt.Errorf("siteType query failed: %v", err)
	}
	response.Stype = siteType

	epoch, err := db.getSiteEpoch(params)
	if err != nil {
		return nil, fmt.Errorf("epoch query failed: %v", err)
	}
	response.Epoch = epoch

	response.ExcCount, response.ArtiCount, response.ExcArea, err = db.getSiteExcArtiProps(params)
	if err != nil {
		return nil, fmt.Errorf("epoch query failed: %v", err)
	}

	heritages, err := db.getSiteHeritages(params)
	if err != nil {
		return nil, fmt.Errorf("heritages query failed: %v", err)
	}
	response.Heritages = heritages

	cultures, err := db.getSiteCultures(params)
	if err != nil {
		return nil, fmt.Errorf("cultures query failed: %v", err)
	}
	response.Cultures = cultures

	resCount, err := db.getSiteResCount(params)
	if err != nil {
		return nil, fmt.Errorf("resCount query failed: %v", err)
	}
	response.ResCount = resCount

	return &response, nil
}

/*
 * Site researches
 */

type siteResearch struct {
	ResID     int64  `json:"res_id"`
	ResName   string `json:"res_name"`
	ResYear   int64  `json:"res_year"`
	ResType   string `json:"res_type"`
	SiteName  string `json:"site_name"`
	Culture   string `json:"culture"`
	ExcCount  int64  `json:"exc_count"`
	ArtiCount int64  `json:"art_count"`
}

func (db *DB) QuerySiteResearches(req map[string]interface{}) ([]*siteResearch, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	res, err := db.getSiteResearches(params)
	if err != nil {
		return nil, err
	}

	return res, nil
}

/*
 * Site reports
 */

type siteReport struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Year   int64  `json:"year"`
	Author string `json:"author"`
}

func (db *DB) QuerySiteReports(req map[string]interface{}) ([]*siteReport, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	rep, err := db.getSiteReports(params)
	if err != nil {
		return nil, err
	}

	return rep, nil
}

/*
 *	Site excavations
 */

type siteExcavation struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Area      float64 `json:"area,omitempty"`
	Boss      string  `json:"boss,omitempty"`
	ResAuthor string  `json:"res_author"`
	ResYear   int64   `json:"res_year"`
}

func (db *DB) QuerySiteExcavations(req map[string]interface{}) ([]*siteExcavation, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	exc, err := db.getSiteExcavations(params)
	if err != nil {
		return nil, err
	}

	return exc, nil
}

/*

 */

type siteArtifact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ResAuthor string `json:"res_author"`
	ResYear   int64  `json:"res_year"`
}

func (db *DB) QuerySiteArtifacts(req map[string]interface{}) ([]*siteArtifact, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	artifacts, err := db.getSiteArtifacts(params)
	if err != nil {
		return nil, err
	}

	return artifacts, nil
}

/*
 * Plural
 */

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

// func (db *DB) Sites(req gin.H) ([]pluralSite, error) {
// 	rows, err := db.QueryNeo(
// 		fmt.Sprintf(pluralstatement, siteFilterString(req)),
// 		gin.H{
// 			"name":   buildRegexpFilter(req["name"]),
// 			"epoch":  req["epoch"],
// 			"type":   req["type"],
// 			"offset": req["offset"],
// 			"limit":  req["limit"],
// 		},
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not query in neo4j: %v", err)
// 	}

// 	var sites []pluralSite
// 	var ids []int64
// 	for err == nil {
// 		var row []interface{}
// 		row, _, err = rows.NextNeo()
// 		if err != nil && err != io.EOF {
// 			return nil, fmt.Errorf("didn't get rows: %v", err)
// 		} else if err == io.EOF {
// 			continue
// 		}

// 		sites = append(sites, pluralSite{
// 			ID: row[0].(int64),
// 			Item: siteItem{
// 				Names:    row[1],
// 				ResNames: row[2],
// 				Epoch:    row[3].(int64),
// 				SiteType: row[4].(int64),
// 			},
// 		})
// 		ids = append(ids, row[0].(int64))
// 	}
// 	err = rows.Close()
// 	if err != nil {
// 		return nil, fmt.Errorf("error when closing statement: %v", err)
// 	}

// 	preStmt := actualCoordinates(len(ids), monument)
// 	coordRows, err := db.QueryNeo(preStmt, idsForQueriengCoordinates(ids))
// 	if err != nil {
// 		return nil, fmt.Errorf("could not query in neo4j: %v", err)
// 	}
// 	defer func() {
// 		err = coordRows.Close()
// 		if err != nil {
// 			fmt.Printf("closing coordinates statement failed: %v", err)
// 		}
// 	}()

// 	coords, _, err := coordRows.All()
// 	if err != nil {
// 		return nil, fmt.Errorf("could not get rows of data: %v", err)
// 	}

// 	for i := range sites {
// 		sites[i].Coords.X = coords[i][0]
// 		sites[i].Coords.Y = coords[i][1]
// 	}

// 	return sites, nil
// }

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
