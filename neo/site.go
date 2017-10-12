package neo

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//////////////////////
////// Singular //////
//////////////////////

type (
	singleSite struct {
		Name     string      `json:"name"`
		Names    interface{} `json:"names"`
		Epoch    string      `json:"epoch"`
		Stype    string      `json:"type"`
		Cultures []string    `json:"cultures"`

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

	siteSpatialReferences struct {
		Date     int64   `json:"date"`
		Accuracy int64   `json:"accuracy"`
		X        float64 `json:"x"`
		Y        float64 `json:"y"`
	}
)

const (
	coordPoint = iota + 1
	coordPolygon
)

func (db *DB) GetSite(id int64, lang string) (interface{}, error) {
	params := gin.H{
		"id":   id,
		"lang": lang,
	}
	queries := []string{
		`MATCH (:Monument {id: {id}})<--(k:Knowledge)
		RETURN COLLECT(k.monument_name) as names`,

		`MATCH (:Monument {id: {id}})-->(sp:SpatialReference)-->(spt:SpatialReferenceType)
		WITH sp, spt
		ORDER BY spt.id ASC, sp.date DESC
		RETURN
			sp.x as x,
			sp.y as y,
			spt.id as accuracy,
			sp.date as date`,

		`MATCH (:Monument {id: {id}})-->(:MonumentType)-[:translation {lang: {lang}}]->(tr:Translate)
		RETURN tr.name as name`,

		`MATCH (:Monument {id: {id}})-->(:Epoch)-[:translation {lang: {lang}}]->(tr:Translate)
		RETURN tr.name as name`,

		`MATCH (:Monument {id: {id}})-->(e:Excavation)
		OPTIONAL MATCH (e)-->(a:Artifact)
		RETURN
			COUNT(e) as excLength,
			SUM(e.area) as excArea,
			COUNT(a) as artiLength`,

		`MATCH (:Monument {id: {id}})<--(n:Heritage)
		RETURN
			n.id as id,
			n.name as name`,

		`MATCH (:Monument {id: {id}})<--(:Knowledge)-->(:Culture)-[:translation {lang: {lang}}]->(tr:Translate)
		RETURN COLLECT(tr.name) as names`,

		`MATCH (:Monument {id: {id}})<--(:Knowledge)<--(r:Research)
		RETURN COUNT(r) as count`,
	}

	pipeline, err := db.PreparePipeline(queries...)
	if err != nil {
		return nil, fmt.Errorf("preparing pipeline failed: %v", err)
	}

	pipeRows, err := pipeline.QueryPipeline(params, params, params, params, params, params, params, params)
	if err != nil {
		return nil, fmt.Errorf("query pipeline failed: %v", err)
	}

	names, _, _, err := pipeRows.NextPipeline()
	if err != nil {
		return nil, fmt.Errorf("couldn't get row from pipeline: %v", err)
	}

	var resp singleSite
	resp.Names = names[0]

	_, _, pipeRows, err = pipeRows.NextPipeline()
	if err != nil {
		return nil, fmt.Errorf("couldn't get end of the row from pipeline: %v", err)
	}

	for {
		sp, _, ref, err := pipeRows.NextPipeline()
		if err != nil {
			return nil, fmt.Errorf("couldn't get row from pipeline: %v", err)
		}
		log.Printf("coord: %#v | ref: %#v", sp, ref)
		if ref == nil {
			resp.Coords = append(resp.Coords, siteSpatialReferences{
				X:        sp[0].(float64),
				Y:        sp[1].(float64),
				Accuracy: sp[2].(int64),
				Date:     sp[3].(int64),
			})
		} else {
			pipeRows = ref
			break
		}
	}

	siteType, _, ref, err := pipeRows.NextPipeline()
	if err != nil {
		return nil, fmt.Errorf("couldn't get row from pipeline: %v", err)
	}
	log.Printf("siteType: %#v | ref: %#v", siteType, ref)

	resp.Stype = siteType[0].(string)
	_, _, pipeRows, err = pipeRows.NextPipeline()
	if err != nil {
		return nil, fmt.Errorf("couldn't get end of the row from pipeline: %v", err)
	}

	err = pipeline.Close()
	if err != nil {
		return nil, fmt.Errorf("error when closing pipeline: %v", err)
	}

	return resp, nil
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
