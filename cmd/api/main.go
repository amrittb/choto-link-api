package main

import (
	"net/http"

	"log"

	"github.com/amrittb/choto-link-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type CreateChotoReq struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

func main() {
	// Load .env file
	godotenv.Load()

	// Initialize etcd
	etcdClient, err := getEtcdClient()
	if err != nil {
		panic(err)
	}
	defer etcdClient.Close()

	// Initialize postgres
	pgxPool, err := getPgxPool()
	if err != nil {
		panic(err)
	}
	defer pgxPool.Close()

	// Build Dependencies
	idAllocator, err := getIdAllocator(etcdClient)
	if err != nil {
		panic(err)
	}
	urlRepository, err := getUrlRepository(pgxPool)
	if err != nil {
		panic(err)
	}
	urlRepository.Initialize()

	shortUrlService := service.NewShortUrlService(idAllocator, urlRepository)

	r := gin.Default()
	r.POST("/api/v1/choto", func(c *gin.Context) {
		var reqJson CreateChotoReq
		if err := c.ShouldBindBodyWithJSON(&reqJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON deserialization error."})
			return
		}

		shortUrl, err := shortUrlService.AllocateAndSaveShortUrl(reqJson.LongUrl)
		if err != nil {
			log.Fatalf("Short Url allocation error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Short Url allocation error."})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"shortUrl": shortUrl,
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

		longUrl, found, err := shortUrlService.GetLongUrl(shortUrl)
		if err != nil {
			log.Fatalf("Short Url fetch error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Short Url fetch error."})
			return
		}

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
