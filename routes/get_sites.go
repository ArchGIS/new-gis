package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	requestParams struct {
		Name   string `form:"name"`
		Epoch  int    `form:"epoch_id"`
		Type   int    `form:"type_id"`
		Offset int    `form:"offset"`
		Limit  int    `form:"limit"`
	}
)

// Sites gets info about archeological sites
func Sites(c *gin.Context) {
	// req := requestParams{Limit: 20}

	// if err := c.Bind(&req); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	// }

	// sites, err := db.Sites(gin.H{
	// 	"name":   req.Name,
	// 	"epoch":  req.Epoch,
	// 	"type":   req.Type,
	// 	"offset": req.Offset,
	// 	"limit":  req.Limit,
	// })
	// if err != nil {
	// 	log.Print(err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"sites": sites})
}

type siteInfoRequest struct {
	Lang string `form:"lang" binding:"eq=en|eq=ru"`
}

// SingleSite get general info about single archaelogical site
func SingleSite(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panicf("could not convert id to int: %v", err)
	}

	req := siteInfoRequest{Lang: "en"}
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panicf("could not bind query params: %v", err)
	}

	site, err := db.GetSite(map[string]interface{}{
		"id":   id,
		"lang": req.Lang,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Panicf("error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"site": site})
}

// SiteResearches get researches related to site
func SiteResearches(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panicf("could not convert id to int: %v", err)
	}

	req := siteInfoRequest{Lang: "en"}
	if err := c.BindQuery(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panicf("could not bind query params: %v", err)
	}

	res, err := db.QuerySiteResearches(map[string]interface{}{
		"id":   id,
		"lang": req.Lang,
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Panicf("query failed: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"site_researches": res})
}

// SiteReports get reports related to site
func SiteReports(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Panicf("could not convert id to int: %v", err)
	}

	reports, err := db.QuerySiteReports(map[string]interface{}{"id": id})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Panicf("query failed: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"site_reports": reports})
}
