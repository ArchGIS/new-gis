package neo

import (
	"fmt"
	"strconv"

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

// GetSite returns general info about site from db
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

// QuerySiteResearches returns site researches from db
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

type siteReport struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Year   int64  `json:"year"`
	Author string `json:"author"`
}

// QuerySiteReports returns site reports from db
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

type siteExcavation struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Area      float64 `json:"area,omitempty"`
	Boss      string  `json:"boss,omitempty"`
	ResAuthor string  `json:"res_author"`
	ResYear   int64   `json:"res_year"`
}

// QuerySiteExcavations returns site excavations from db
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
 * Site artifacts
 */

type siteArtifact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ResAuthor string `json:"res_author"`
	ResYear   int64  `json:"res_year"`
}

// QuerySiteArtifacts returns site artifacts from db
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

type siteCarbon struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Date     int64  `json:"date"`
	Sigma    int64  `json:"sigma"`
	Material string `json:"material"`
}

// QuerySiteRadioCarbon returns site radiocarbon from db
func (db *DB) QuerySiteRadioCarbon(req map[string]interface{}) ([]*siteCarbon, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	radiocarbon, err := db.getSiteRadioCarbon(params)
	if err != nil {
		return nil, err
	}

	return radiocarbon, nil
}

type sitePhoto struct {
	ID        string `json:"fileid"`
	Part      string `json:"part"`
	Direction string `json:"direction"`
}

func (db *DB) QuerySitePhotos(req map[string]interface{}) ([]*sitePhoto, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	photos, err := db.getSitePhotos(params)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

type siteTopo struct {
	ID     string `json:"fileid"`
	Author string `json:"author"`
	Year   int64  `json:"year"`
}

func (db *DB) QuerySiteTopoplans(req map[string]interface{}) ([]*siteTopo, error) {
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	topos, err := db.getSiteTopos(params)
	if err != nil {
		return nil, err
	}

	return topos, nil
}

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

func (db *DB) Sites(req map[string]interface{}) ([]*site, error) {
	sites, err := db.getSites(req)
	if err != nil {
		return nil, err
	}

	return sites, nil
}

func idsForQueriengCoordinates(ids []int64) gin.H {
	params := make(gin.H)
	for i, v := range ids {
		params["id"+strconv.Itoa(i)] = v
	}

	return params
}
