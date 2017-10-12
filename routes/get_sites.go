package routes

import (
	"log"
	"net/http"

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
	req := requestParams{Limit: 20}

	if err := c.Bind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	}

	sites, err := db.Sites(gin.H{
		"name":   req.Name,
		"epoch":  req.Epoch,
		"type":   req.Type,
		"offset": req.Offset,
		"limit":  req.Limit,
	})
	if err != nil {
		log.Print(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"sites": sites})
}

// SingleSite get info about single archaelogical site
func SingleSite(c *gin.Context) {
	id := c.Param("id")
	lang := c.Query("lang")
	if lang == "" {
		lang = "en"
	}

	site, err := db.GetSite(id, lang)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"site": site})
}

// SiteResearches get researches related to site
func SiteResearches(c *gin.Context) {
	id := c.Param("id")
	lang := c.Query("lang")
	if lang == "" {
		lang = "en"
	}

	res, err := db.QuerySiteResearches(id, lang)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{"site_researches": res})
}
