package routes

import (
  "io/ioutil"
  "log"
  "encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
  "strings"
  "strconv"
)

type Parameter struct {
  Name string;
  Value string;
  Not bool;
}
type SearchQuery struct {
  Id  int
  Params []Parameter
}

type QueryTemp struct {
  Match string
  Where string
}

func Example(c *gin.Context) {
  req, err := ioutil.ReadAll(c.Request.Body)
  if err != nil {
    log.Panicf("cannot read request body: %v", err)
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  var SQuery []SearchQuery
  err = json.Unmarshal(req, &SQuery)
  if err != nil {
    log.Panicf("cannot read request body: %v", err)
    c.AbortWithStatus(http.StatusBadRequest)
    return
  }

  log.Print("DATA: %v", SQuery) 
  result := SendQuery(SQuery, c);
	c.JSON(http.StatusOK, result)
}

func SendQuery(params []SearchQuery, c *gin.Context) interface{} {
  templates := GetTemplates()
  queryStr := ""

  for t, query := range params {

    if len(query.Params) == 0 {
      continue
    }

    Match := ""
    Where := ""
    Ret := ""  
    entity := ""

    if query.Id == 1 {
      entity = "site-"
    }
    if query.Id == 2 {
      entity = "res-"
    }
    if query.Id == 3 {
      entity = "auth-"
    }

    Match = templates[entity+"main"].Match
    Ret = templates[entity+"main"].Where

    for i, param := range query.Params {
      count := strconv.Itoa(i)
      m := strings.Replace(templates[entity+param.Name].Match + " ", "_ID", count, -1)
      w := strings.Replace(templates[entity+param.Name].Where + " ", "_ID", count, -1)
      w = strings.Replace(w, "_VALUE", param.Value, -1)
      
      if param.Not == true {
        w = "NOT " + w
      }

      Match += m
      Where += w
      
      if i+1 < len(query.Params) {
        Where += " AND "
      }
    }

    if t != 0 {
        queryStr += " UNION "
    }
    queryStr += Match + " WHERE " + Where + Ret
    log.Print("QUERY %d: %v", t, queryStr) 

  }

  results, err := db.RawQuery(queryStr)
  if err != nil {
    log.Panicf("cannot read request body: %v", err)
    c.AbortWithStatus(http.StatusBadRequest)
  }
  log.Print("result: %v", results) 

  return results;
}

func GetTemplates() map[string]QueryTemp {
  qt := make(map[string]QueryTemp)
  qt["auth-main"] = QueryTemp{ "MATCH (a:Author) ", "RETURN DISTINCT a"}

  qt["auth-authName"] = QueryTemp{ "", "a.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-authJob"] = QueryTemp{ "MATCH (a)--(o_ID:Organization)--(j_ID:AuthorJob)", "j_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  qt["auth-artCul"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)--(c_ID:Culture)", "c_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-artName"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)", "art_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-artMat"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)--(artmat_ID:ArtifactMaterial)", "artmat_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-artCat"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)--(artcat_ID:ArtifactCategory)", "artcat_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-artYearBefore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)", "art_ID.year <= _VALUE"}
  qt["auth-artYearAfter"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)", "art_ID.year >= _VALUE"}
  qt["auth-artPhotoBefore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)--(i_ID:Image)", "i_ID.year <= _VALUE"}
  qt["auth-artPhotoAfter"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)--(i_ID:Image)", "i_ID.year >= _VALUE"}

  qt["auth-excAreaMore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(e_ID:Excavation)", "e_ID.area >= _VALUE"}
  qt["auth-excAreaLess"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(e_ID:Excavation)", "e_ID.area <= _VALUE"}
  qt["auth-excObjName"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(e_ID:Excavation)--(obj_ID:Complex)", "obj_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-excBossName"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(e_ID:Excavation)", "e_ID.boss =~ '(?ui)^.*(_VALUE).*$'"}

  qt["auth-resAfter"] = QueryTemp{ "MATCH (a)--(r_ID:Research)", "r_ID.year >= _VALUE"}
  qt["auth-resBefore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)", "r_ID.year <= _VALUE"}
  qt["auth-resType"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(type_ID:ResearchType)", "type_ID.id = _VALUE"}

  qt["auth-colStorage"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact)--(int_ID:StorageInterval)--(coll_ID:Collection)--(org_ID:Organization)", "org_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-colName"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(art_ID:Artifact--(int_ID:StorageInterval)--(coll_ID:Collection)", "coll_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  qt["auth-herSec"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)--(type_ID:SecurityType)", "type_ID.id = _VALUE"}
  qt["auth-herReg"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)", "her_ID.code =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-herAfter"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)", "her_ID.docDate >= _VALUE"}
  qt["auth-herBefore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)", "her_ID.docDate <= _VALUE"}

  qt["auth-siteName"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)", "ks_ID.monument_name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-siteCul"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)--(cult_ID:Culture)", "cult_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-siteType"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(type_ID:MonumentType)", "type_ID.id = _VALUE"}
  qt["auth-siteEpoch"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(epoch_ID:Epoch)", "epoch_ID.id = _VALUE"}
  qt["auth-siteTopoAfter"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:hastopo]-(topo_ID:Image)", "photo_ID.creationDate >= _VALUE"}
  qt["auth-siteTopoBefore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:hastopo]-(topo_ID:Image)", "photo_ID.creationDate <= _VALUE"}
  qt["auth-sitePhotoAfter"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:has]-(photo_ID:Image)", "photo_ID.creationDate >= _VALUE"}
  qt["auth-sitePhotoBefore"] = QueryTemp{ "MATCH (a)--(r_ID:Research)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:has]-(photo_ID:Image)", "photo_ID.creationDate <= _VALUE"}

  qt["auth-pubName"] = QueryTemp{ "MATCH (a)--(pub_ID:Publication)", "pub_ID.publication_name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-pubPlace"] = QueryTemp{ "MATCH (a)--(pub_ID:Publication)", "pub_ID.publicated_in =~ '(?ui)^.*(_VALUE).*$'"}
  qt["auth-pubTitle"] = QueryTemp{ "MATCH (a)--(pub_ID:Publication)", "pub_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  ////////////////////////////////////////////////////////////////////////////////////////
  qt["res-main"] = QueryTemp{ "MATCH (a:Research) ", "RETURN DISTINCT a"}

  qt["res-authName"] = QueryTemp{ "MATCH (a)--(author_ID:Author)", "author_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-authJob"] = QueryTemp{ "MATCH (a)--(author_ID:Author)--(o_ID:Organization)--(j_ID:AuthorJob)", "j_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  qt["res-artCul"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)--(c_ID:Culture)", "c_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-artName"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)", "art_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-artMat"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)--(artmat_ID:ArtifactMaterial)", "artmat_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-artCat"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)--(artcat_ID:ArtifactCategory)", "artcat_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-artYearBefore"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)", "art_ID.year <= _VALUE"}
  qt["res-artYearAfter"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)", "art_ID.year >= _VALUE"}
  qt["res-artPhotoBefore"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)--(i_ID:Image)", "i_ID.year <= _VALUE"}
  qt["res-artPhotoAfter"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)--(i_ID:Image)", "i_ID.year >= _VALUE"}

  qt["res-excAreaMore"] = QueryTemp{ "MATCH (a)--(e_ID:Excavation)", "e_ID.area >= _VALUE"}
  qt["res-excAreaLess"] = QueryTemp{ "MATCH (a)--(e_ID:Excavation)", "e_ID.area <= _VALUE"}
  qt["res-excObjName"] = QueryTemp{ "MATCH (a)--(e_ID:Excavation)--(obj_ID:Complex)", "obj_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-excBossName"] = QueryTemp{ "MATCH (a)--(e_ID:Excavation)", "e_ID.boss =~ '(?ui)^.*(_VALUE).*$'"}

  qt["res-resAfter"] = QueryTemp{ "", "a.year >= _VALUE"}
  qt["res-resBefore"] = QueryTemp{ "", "a.year <= _VALUE"}
  qt["res-resType"] = QueryTemp{ "MATCH (a)--(type_ID:ResearchType)", "type_ID.id = _VALUE"}

  qt["res-colStorage"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact)--(int_ID:StorageInterval)--(coll_ID:Collection)--(org_ID:Organization)", "org_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-colName"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(art_ID:Artifact--(int_ID:StorageInterval)--(coll_ID:Collection)", "coll_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  qt["res-herSec"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)--(type_ID:SecurityType)", "type_ID.id = _VALUE"}
  qt["res-herReg"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)", "her_ID.code =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-herAfter"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)", "her_ID.docDate >= _VALUE"}
  qt["res-herBefore"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(her_ID:Heritage)", "her_ID.docDate <= _VALUE"}

  qt["res-siteName"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)", "ks_ID.monument_name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-siteCul"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)--(cult_ID:Culture)", "cult_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-siteType"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(type_ID:MonumentType)", "type_ID.id = _VALUE"}
  qt["res-siteEpoch"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(epoch_ID:Epoch)", "epoch_ID.id = _VALUE"}
  qt["res-siteTopoAfter"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:hastopo]-(topo_ID:Image)", "photo_ID.creationDate >= _VALUE"}
  qt["res-siteTopoBefore"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:hastopo]-(topo_ID:Image)", "photo_ID.creationDate <= _VALUE"}
  qt["res-sitePhotoAfter"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:has]-(photo_ID:Image)", "photo_ID.creationDate >= _VALUE"}
  qt["res-sitePhotoBefore"] = QueryTemp{ "MATCH (a)--(k_ID:Knowledge)--(site_ID:Monument)--(ks_ID:Knowledge)-[:has]-(photo_ID:Image)", "photo_ID.creationDate <= _VALUE"}

  qt["res-pubName"] = QueryTemp{ "MATCH (a)--(author_ID:Author)--(pub_ID:Publication)", "pub_ID.publication_name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-pubPlace"] = QueryTemp{ "MATCH (a)--(author_ID:Author)--(pub_ID:Publication)", "pub_ID.publicated_in =~ '(?ui)^.*(_VALUE).*$'"}
  qt["res-pubTitle"] = QueryTemp{ "MATCH (a)--(author_ID:Author)--(pub_ID:Publication)", "pub_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  /////////////////////////////////////////////////////////////////////////////
  qt["site-main"] = QueryTemp{ "MATCH (mt:MonumentType)--(a:Monument)--(k:Knowledge) MATCH (a)--(e:Epoch)", "RETURN DISTINCT {id: a.id, name: k.monument_name, epoch: e.name, type: mt.name}"}

  qt["site-authName"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(author_ID:Author)", "author_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-authJob"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(author_ID:Author)--(o_ID:Organization)--(j_ID:AuthorJob)", "j_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  qt["site-artCul"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)--(c_ID:Culture)", "c_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-artName"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)", "art_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-artMat"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)--(artmat_ID:ArtifactMaterial)", "artmat_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-artCat"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)--(artcat_ID:ArtifactCategory)", "artcat_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-artYearBefore"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)", "art_ID.year <= _VALUE"}
  qt["site-artYearAfter"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)", "art_ID.year >= _VALUE"}
  qt["site-artPhotoBefore"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)--(i_ID:Image)", "i_ID.year <= _VALUE"}
  qt["site-artPhotoAfter"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)--(i_ID:Image)", "i_ID.year >= _VALUE"}

  qt["site-excAreaMore"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(e_ID:Excavation)", "e_ID.area >= _VALUE"}
  qt["site-excAreaLess"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(e_ID:Excavation)", "e_ID.area <= _VALUE"}
  qt["site-excObjName"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(e_ID:Excavation)--(obj_ID:Complex)", "obj_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-excBossName"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(e_ID:Excavation)", "e_ID.boss =~ '(?ui)^.*(_VALUE).*$'"}

  qt["site-resAfter"] = QueryTemp{ "MATCH (k)--(res_ID:Research)", "res_ID.year >= _VALUE"}
  qt["site-resBefore"] = QueryTemp{ "MATCH (k)--(res_ID:Research)", "res_ID.year <= _VALUE"}
  qt["site-resType"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(type_ID:ResearchType)", "type_ID.id = _VALUE"}

  qt["site-colStorage"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact)--(int_ID:StorageInterval)--(coll_ID:Collection)--(org_ID:Organization)", "org_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-colName"] = QueryTemp{ "MATCH (k)--(art_ID:Artifact--(int_ID:StorageInterval)--(coll_ID:Collection)", "coll_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  qt["site-herSec"] = QueryTemp{ "MATCH (a)--(her_ID:Heritage)--(type_ID:SecurityType)", "type_ID.id = _VALUE"}
  qt["site-herReg"] = QueryTemp{ "MATCH (a)--(her_ID:Heritage)", "her_ID.code =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-herAfter"] = QueryTemp{ "MATCH (a)--(her_ID:Heritage)", "her_ID.docDate >= _VALUE"}
  qt["site-herBefore"] = QueryTemp{ "MATCH (a)--(her_ID:Heritage)", "her_ID.docDate <= _VALUE"}

  qt["site-siteName"] = QueryTemp{ "MATCH (a)--(ks_ID:Knowledge)", "ks_ID.monument_name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-siteCul"] = QueryTemp{ "MATCH (a)--(ks_ID:Knowledge)--(cult_ID:Culture)", "cult_ID.name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-siteType"] = QueryTemp{ "MATCH (a)--(type_ID:MonumentType)", "type_ID.id = _VALUE"}
  qt["site-siteEpoch"] = QueryTemp{ "MATCH (a)--(epoch_ID:Epoch)", "epoch_ID.id = _VALUE"}
  qt["site-siteTopoAfter"] = QueryTemp{ "MATCH (a)--(ks_ID:Knowledge)-[:hastopo]-(topo_ID:Image)", "photo_ID.creationDate >= _VALUE"}
  qt["site-siteTopoBefore"] = QueryTemp{ "MATCH (a)--(ks_ID:Knowledge)-[:hastopo]-(topo_ID:Image)", "photo_ID.creationDate <= _VALUE"}
  qt["site-sitePhotoAfter"] = QueryTemp{ "MATCH (a)--(ks_ID:Knowledge)-[:has]-(photo_ID:Image)", "photo_ID.creationDate >= _VALUE"}
  qt["site-sitePhotoBefore"] = QueryTemp{ "MATCH (a)--(ks_ID:Knowledge)-[:has]-(photo_ID:Image)", "photo_ID.creationDate <= _VALUE"}

  qt["site-pubName"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(author_ID:Author)--(pub_ID:Publication)", "pub_ID.publication_name =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-pubPlace"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(author_ID:Author)--(pub_ID:Publication)", "pub_ID.publicated_in =~ '(?ui)^.*(_VALUE).*$'"}
  qt["site-pubTitle"] = QueryTemp{ "MATCH (k)--(res_ID:Research)--(author_ID:Author)--(pub_ID:Publication)", "pub_ID.name =~ '(?ui)^.*(_VALUE).*$'"}

  return qt 
}
