package router

import (
	"fmt"
	"github.com/ankurgel/reducto/internal/store"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Router struct {
	Engine *gin.Engine
}


func InitRouter(s *store.Store) *Router {
	r := gin.Default()
	r.Use(dbMiddleware(s))

	r.GET("/", rootHandler)
	r.GET("/:shortUrl", longV1Handler)
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{
		v1.POST("/shorten", shortenV1Handler)
	}


	return &Router{r}
}

func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func shortenV1Handler(c *gin.Context) {
	s := c.MustGet("store").(*store.Store)
	longUrl := c.PostForm("url")
	customSlugRequested := c.PostForm("custom")
	result, e := s.CreateByLongURL(longUrl, customSlugRequested)
	if e != nil {
		errorMessage := fmt.Sprintf("Error in shortening %s : %s", longUrl, e)
		log.Error(errorMessage)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errorMessage})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"shortUrl": result.ShortURL(), "longUrl": result.Original})
}

func longV1Handler(c *gin.Context) {
	s := c.MustGet("store").(*store.Store)
	shortUrl := c.Param("shortUrl")
	url, err := s.FindByShortURL(shortUrl)
	if err != nil {
		errorMessage := fmt.Sprintf("Error in getSlug for %s: %s", shortUrl, err)
		log.Error(errorMessage)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url.Original)
}

func dbMiddleware(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("store", s)
		c.Next()
	}
}