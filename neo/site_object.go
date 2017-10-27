package neo

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/johnnadratowski/golang-neo4j-bolt-driver/encoding"
)

func (db *DB) getSiteNames(params []byte) ([]string, error) {
	statement := `MATCH (:Monument {id: {id}})<--(k:Knowledge)
		RETURN k.monument_name`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		names = append(names, name)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return names, nil
}

func (db *DB) getSiteCoordinates(params []byte) ([]*siteSpatialReferences, error) {
	statement := `
		MATCH (:Monument {id: {id}})-->(sp:SpatialReference)-->(spt:SpatialReferenceType)
		WITH sp, spt
		ORDER BY spt.id ASC, sp.date DESC
		RETURN
			sp.x as x,
			sp.y as y,
			spt.id as accuracy,
			sp.date as date
	`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coords []*siteSpatialReferences
	for rows.Next() {
		coord := new(siteSpatialReferences)
		err = rows.Scan(&coord.X, &coord.Y, &coord.Accuracy, &coord.Date)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		coords = append(coords, coord)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return coords, nil
}

func (db *DB) getSiteType(params []byte) (string, error) {
	statement := `MATCH (:Monument {id: {id}})-->(:MonumentType)-[:translation {lang: {lang}}]->(tr:Translate)
		RETURN tr.name as name`

	var stype string
	err := db.QueryRow(statement, params).Scan(&stype)
	if err != nil {
		return "", err
	}

	return stype, nil
}

func (db *DB) getSiteEpoch(params []byte) (string, error) {
	statement := `MATCH (:Monument {id: {id}})-->(:Epoch)-[:translation {lang: {lang}}]->(tr:Translate)
		RETURN tr.name as name`

	var epoch string
	err := db.QueryRow(statement, params).Scan(&epoch)
	if err != nil {
		return "", err
	}

	return epoch, nil
}

func (db *DB) getSiteExcArtiProps(params []byte) (int64, int64, float64, error) {
	statement := `MATCH (:Monument {id: {id}})-->(e:Excavation)
		OPTIONAL MATCH (e)-->(a:Artifact)
		RETURN
			COUNT(e) as excLength,
			SUM(e.area) as excArea,
			COUNT(a) as artiLength`

	var excLen, artiLen int64
	var excArea float64
	err := db.QueryRow(statement, params).Scan(&excLen, &excArea, &artiLen)
	if err != nil {
		return 0, 0, 0, err
	}

	return excLen, artiLen, excArea, nil
}

func (db *DB) getSiteHeritages(params []byte) ([]*nHeritage, error) {
	statement := `MATCH (:Monument {id: {id}})<--(n:Heritage)
		RETURN
			n.id as id,
			n.name as name`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var heritages []*nHeritage
	for rows.Next() {
		herit := new(nHeritage)
		err = rows.Scan(&herit.ID, &herit.Name)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		heritages = append(heritages, herit)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return heritages, nil
}

func (db *DB) getSiteCultures(params []byte) ([]string, error) {
	statement := `MATCH (:Monument {id: {id}})<--(:Knowledge)-->(:Culture)-[:translation {lang: {lang}}]->(tr:Translate)
		RETURN tr.name`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cultures []string
	for rows.Next() {
		var culture string
		err = rows.Scan(&culture)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		cultures = append(cultures, culture)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return cultures, nil
}

func (db *DB) getSiteResCount(params []byte) (int64, error) {
	statement := `MATCH (:Monument {id: {id}})<--(:Knowledge)<--(r:Research)
		RETURN COUNT(r) as count`

	var resCount int64
	err := db.QueryRow(statement, params).Scan(&resCount)
	if err != nil {
		return 0, err
	}

	return resCount, nil
}

func (db *DB) getSiteResearches(params []byte) ([]*siteResearch, error) {
	statement := `MATCH (s:Monument {id: {id}})
		MATCH (s)<--(k:Knowledge)<--(r:Research)-->(rt:ResearchType)-[:translation {lang: {lang}}]->(rtTr:Translate)
		MATCH (k)-->(c:Culture)-[:translation {lang: {lang}}]->(cTr:Translate)
		OPTIONAL MATCH (s)-->(e:Excavation)<--(r)
		OPTIONAL MATCH (e)-->(a:Artifact)
		return r.id, r.name, r.year, rtTr.name, k.monument_name, cTr.name, COUNT(e), COUNT(a)`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var researches []*siteResearch
	for rows.Next() {
		res := new(siteResearch)
		err = rows.Scan(&res.ResID, &res.ResName, &res.ResYear, &res.ResType, &res.SiteName, &res.Culture, &res.ExcCount, &res.ArtiCount)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		researches = append(researches, res)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return researches, nil
}

func (db *DB) getSiteReports(params []byte) ([]*siteReport, error) {
	statement := `MATCH (:Monument {id: {id}})<--(:Knowledge)<--(:Research)-->(rep:Report)-->(a:Author)
		RETURN rep.id, rep.name, rep.year, a.name`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*siteReport
	for rows.Next() {
		rep := new(siteReport)
		err = rows.Scan(&rep.ID, &rep.Name, &rep.Year, &rep.Author)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}
		reports = append(reports, rep)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return reports, nil
}

func (db *DB) getSiteExcavations(params []byte) ([]*siteExcavation, error) {
	statement := `MATCH (s:Monument {id: {id}})<--(:Knowledge)<--(r:Research)-[:hasauthor]->(a:Author)
	MATCH (s)-->(e:Excavation)<--(r)
	RETURN e.id, e.name, e.area, e.boss, a.name, r.year`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var excs []*siteExcavation
	for rows.Next() {
		var boss sql.NullString
		var area sql.NullFloat64
		exc := new(siteExcavation)
		err = rows.Scan(&exc.ID, &exc.Name, &area, &boss, &exc.ResAuthor, &exc.ResYear)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}

		if boss.Valid {
			exc.Boss = boss.String
		}
		if area.Valid {
			exc.Area = area.Float64
		}

		excs = append(excs, exc)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return excs, nil
}

func (db *DB) getSiteArtifacts(params []byte) ([]*siteArtifact, error) {
	statement := `MATCH (s:Monument {id: {id}})<--(:Knowledge)<--(r:Research)-[:hasauthor]->(aut:Author)
	MATCH (s)-->(e:Excavation)<--(r)
	MATCH (e)-->(art:Artifact)
	RETURN art.id, art.name, aut.name, r.year`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var artifacts []*siteArtifact
	for rows.Next() {
		art := new(siteArtifact)
		err = rows.Scan(&art.ID, &art.Name, &art.ResAuthor, &art.ResYear)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}

		artifacts = append(artifacts, art)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return artifacts, nil
}

func (db *DB) getSiteRadioCarbon(params []byte) ([]*siteCarbon, error) {
	statement := `MATCH (:Monument {id: {id}})<--(:Knowledge)-->(rc:Radiocarbon)-->(cm:CarbonMaterial)
	RETURN rc.id, rc.name, rc.date, rc.s, cm.name`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carbon []*siteCarbon
	for rows.Next() {
		rc := new(siteCarbon)
		err = rows.Scan(&rc.ID, &rc.Name, &rc.Date, &rc.Sigma, &rc.Material)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}

		carbon = append(carbon, rc)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return carbon, nil
}

func (db *DB) getSitePhotos(params []byte) ([]*sitePhoto, error) {
	statement := `MATCH (:Monument {id: {id}})<--(:Knowledge)-[:has]->(i:Image)-->(cd:CardinalDirection)
	RETURN i.fileid, i.monumentPart, cd.name`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*sitePhoto
	for rows.Next() {
		photo := new(sitePhoto)
		err = rows.Scan(&photo.ID, &photo.Part, &photo.Direction)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}

		photos = append(photos, photo)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return photos, nil
}

func (db *DB) getSiteTopos(params []byte) ([]*siteTopo, error) {
	statement := `MATCH (:Monument {id: {id}})<--(:Knowledge)-[:hastopo]->(i:Image)
	RETURN i.fileid, i.author, i.year`

	rows, err := db.Query(statement, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topos []*siteTopo
	for rows.Next() {
		topo := new(siteTopo)
		err = rows.Scan(&topo.ID, &topo.Author, &topo.Year)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}

		topos = append(topos, topo)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return topos, nil
}

type site struct {
	ID       int64    `json:"id"`
	Names    []string `json:"site_names"`
	ResNames []string `json:"res_names"`
	Epoch    int64    `json:"epoch_id"`
	Type     int64    `json:"type_id"`
}

func (db *DB) getSites(req map[string]interface{}) ([]*site, error) {
	stmt := `
    MATCH (s:Monument)<--(k:Knowledge)
    MATCH (s)-[:has]->(st:MonumentType)
    MATCH (s)-[:has]->(e:Epoch)
		%s
		WITH
			s.id AS id,
			e.id AS epoch,
			st.id AS type
    RETURN DISTINCT id, epoch, type
    SKIP {offset} LIMIT {limit}
	`
	stmtWithFilters := fmt.Sprintf(stmt, siteFilterString(req))
	addRegexpFilter(req, "name")
	params, err := encoding.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %v", err)
	}

	rows, err := db.Query(stmtWithFilters, params)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []*site
	for rows.Next() {
		site := new(site)
		err = rows.Scan(&site.ID, &site.Epoch, &site.Type)
		if err != nil {
			return nil, fmt.Errorf("iterating rows failed: %v", err)
		}

		// Site names and researches names section
		namesParams, err := encoding.Marshal(map[string]interface{}{"id": site.ID})
		if err != nil {
			return nil, err
		}

		namesRows, err := db.Query(`
			MATCH (s:Monument {id: {id}})<--(k:Knowledge)
			MATCH (r:Research)-[:has]->(k)
			RETURN k.monument_name, r.name
		`, namesParams)
		if err != nil {
			return nil, err
		}

		for namesRows.Next() {
			var siteName, resName string
			err = namesRows.Scan(&siteName, &resName)
			if err != nil {
				return nil, err
			}

			site.Names = append(site.Names, siteName)
			site.ResNames = append(site.ResNames, resName)
		}
		if err = rows.Err(); err != nil {
			return nil, fmt.Errorf("end of the rows failed: %v", err)
		}
		err = namesRows.Close()
		if err != nil {
			return nil, err
		}

		sites = append(sites, site)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("end of the rows failed: %v", err)
	}

	return sites, nil
}

func siteFilterString(reqParams map[string]interface{}) string {
	var filter []string
	var stmt string

	if reqParams["name"].(string) != "" {
		filter = append(filter, "k.monument_name =~ {name}")
	}
	if reqParams["epoch"].(int64) != 0 {
		filter = append(filter, "e.id = {epoch}")
	}
	if reqParams["type"].(int64) != 0 {
		log.Printf("type: %#v .", reqParams["type"])
		filter = append(filter, "st.id = {type}")
	}
	if len(filter) > 0 {
		stmt = "WHERE " + strings.Join(filter, " AND ")
	}

	return stmt
}
