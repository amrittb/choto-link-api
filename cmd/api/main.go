package main

import (
	"net/http"

	"github.com/amrittb/choto-link-api/internal/chotoapi"
	"github.com/gin-gonic/gin"
)

type CreateChotoReq struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

func main() {
	r := gin.Default()
	r.POST("/api/v1/choto", func(c *gin.Context) {
		var reqJson CreateChotoReq
		if err := c.ShouldBindBodyWithJSON(&reqJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON deserialization error."})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"shortUrl": chotoapi.Create(reqJson.LongUrl),
		})
	})

	r.GET("/api/v1/choto/:shortUrl", func(c *gin.Context) {
		shortUrl := c.Params.ByName("shortUrl")
		if shortUrl == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "shortUrl not provided.",
			})
			return
		}

		longUrl, found := chotoapi.Get(shortUrl)
		if !found {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Redirection URL not found.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"longUrl": longUrl,
		})
	})

	r.Run(":8081")
}
