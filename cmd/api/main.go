package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amrittb/choto-link-api/internal/allocator"
	"github.com/amrittb/choto-link-api/internal/repository"
	"github.com/amrittb/choto-link-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type CreateChotoReq struct {
	LongUrl string `json:"longUrl" binding:"required"`
}

func getEtcdClient() (*clientv3.Client, error) {
	etcdEndpointRaw := os.Getenv("ETCD_ENDPOINT")
	if etcdEndpointRaw == "" {
		return nil, fmt.Errorf("ETCD_ENDPOINT environment variable is empty.")
	}

	etcdEndpoints := strings.Split(etcdEndpointRaw, ",")
	if len(etcdEndpoints) == 0 {
		return nil, fmt.Errorf("ETCD_ENDPOINT environment variable is empty.")
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func getIdAllocator() (*allocator.IdAllocator, error) {
	etcdClient, err := getEtcdClient()
	if err != nil {
		return nil, err
	}

	etcdKey := os.Getenv("ETCD_ID_ALLOCATION_KEY")
	if etcdKey == "" {
		return nil, fmt.Errorf("ETCD_ID_ALLOCATION_KEY environment variable is empty.")
	}
	rangeSize, err := strconv.ParseUint(os.Getenv("ID_ALLOCATION_RANGE_SIZE"), 10, 64)
	if err != nil {
		return nil, err
	}

	return allocator.NewIdAllocator(etcdClient, etcdKey, rangeSize)
}

func getUrlRepository() (*repository.UrlRepository, error) {
	return repository.NewUrlRepository(), nil
}

func main() {
	// Load .env file
	godotenv.Load()

	// Build Dependencies
	idAllocator, err := getIdAllocator()
	if err != nil {
		panic(err)
	}

	urlRepository, err := getUrlRepository()
	if err != nil {
		panic(err)
	}

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

		longUrl, found := shortUrlService.GetLongUrl(shortUrl)
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
