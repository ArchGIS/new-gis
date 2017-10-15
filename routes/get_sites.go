package routes

import (
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

// SingleSite get info about single archaelogical site
func SingleSite(c *gin.Context) {
	// sID := c.Param("id")
	// ID, err := strconv.Atoi(sID)
	// if err != nil {
	// 	log.Printf("couldn't convert id to integer: %v", err)
	// 	c.AbortWithError(http.StatusNotAcceptable, err)
	// 	return
	// }
	// lang := c.Query("lang")
	// if lang != "" && (lang != "ru" || lang != "en") {
	// 	log.Printf("wrong lang: %v", lang)
	// 	c.AbortWithStatus(http.StatusNotAcceptable)
	// 	return
	// }

	// site, err := db.GetSite(int64(ID), lang)
	// if err != nil {
	// 	log.Printf("error: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"site": site})
}

// SiteResearches get researches related to site
func SiteResearches(c *gin.Context) {
	// id := c.Param("id")
	// lang := c.Query("lang")
	// if lang == "" {
	// 	lang = "en"
	// }

	// res, err := db.QuerySiteResearches(id, lang)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// }

	// c.JSON(http.StatusOK, gin.H{"site_researches": res})
}
