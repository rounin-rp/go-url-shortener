package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rounin-rp/go-url-shortener/mongodb"
	"github.com/rounin-rp/go-url-shortener/shortener"
)

type UrlCreationRequest struct {
	LongUrl string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func CreateShortUrl(c *gin.Context) {
	var creationRequest UrlCreationRequest
	if err := c.ShouldBindJSON(&creationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	shortUrl := shortener.GenerateShortLink(creationRequest.LongUrl, creationRequest.UserId)
	// store.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId)
	mongodb.SaveUrlMapping(shortUrl, creationRequest.LongUrl, creationRequest.UserId)

	host := "http://localhost:9808/"
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": host + shortUrl,
	})
}

func HandleShortUrlRedirect(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	// initialUrl := store.RetrieveInitialUrl(shortUrl)
	initialUrl, err := mongodb.RetrieveInitialUrl(shortUrl)
	if err != nil {
		c.Redirect(404, "www.example.com")
	} else {
		c.Redirect(302, initialUrl)
	}
}
