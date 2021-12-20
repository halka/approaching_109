package main

import (
	"approaching_109/parser"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/approaching/:startId/:poleId", func(c *gin.Context) {
		// CORS
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Authorization")

		params := parser.SetParams(c.Params.ByName("startId"), c.Params.ByName("poleId"))
		c.JSON(http.StatusOK, parser.ParseApproaching(params))
	})

	return r
}
func main() {
	r := setupRouter()
	r.Run(":3001")
}
