package main

import (
	"approaching_109/parser"
	"github.com/gin-gonic/gin"
	"net/http"
)
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/approaching/:startId/:poleId", func(c *gin.Context) {
		params := parser.SetParams(c.Params.ByName("startId"), c.Params.ByName("poleId"))
		c.JSON(http.StatusOK, parser.ParseApproaching(params))
	})

	return r
}
func main() {
	r := setupRouter()
	r.Run(":3001")
}
