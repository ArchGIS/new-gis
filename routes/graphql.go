package routes

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type graphqlRequest struct {
	Query string                 `form:"query" binding:"required"`
	Vars  map[string]interface{} `form:"variables"`
}

// Graphql return ...
func Graphql(c *gin.Context) {
	// var req graphqlRequest
	// if err := c.BindJSON(&req); err != nil {
	// 	c.AbortWithError(http.StatusBadRequest, NotAllowedQueryParams)
	// }
	defer c.Request.Body.Close()
	req, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panicf("cannot read request body: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := db.Graphql(req)
	if err != nil {
		log.Panicf("graphql request failed: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, result)
}
