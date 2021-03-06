package router

import (
	"fmt"
	"github.com/ankurgel/reducto/internal/store"
	"github.com/ankurgel/reducto/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Router struct {
	Engine *gin.Engine
}


func InitRouter(s *store.Store) *Router {
	r := gin.Default()
	r.Use(dbMiddleware(s))

	r.LoadHTMLGlob("templates/*")
	r.GET("/", rootHandler)
	r.GET("/:shortUrl/preview", previewHandler)
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
	sanitizedLongUrl, err := util.NormalizeURL(longUrl, s)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	result, err := s.CreateByLongURL(sanitizedLongUrl, customSlugRequested)
	if err != nil {
		errorMessage := fmt.Sprintf("Error in shortening %s : %s", longUrl, err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": errorMessage})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"shortUrl": result.ShortURL(), "longUrl": result.Original})
}

func longV1Handler(c *gin.Context) {
	s := c.MustGet("store").(*store.Store)
	shortUrl := c.Param("shortUrl")
	if  strings.HasSuffix(shortUrl, "+") {
		previewHandler(c)
		return
	}
	url, err := s.FindByShortURL(shortUrl)
	if err != nil {
		errorMessage := fmt.Sprintf("Error in getSlug for %s: %s", shortUrl, err)
		c.JSON(http.StatusNotFound, gin.H{"error": errorMessage})
		return
	}
	if _, err = s.IncreaseVisitForUrl(url, c.ClientIP()); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.Redirect(http.StatusMovedPermanently, url.Original)
}

func previewHandler(c *gin.Context) {
	s := c.MustGet("store").(*store.Store)
	shortUrl := c.Param("shortUrl")
	url, err := s.FindByShortURL(shortUrl)
	if err != nil {
		errorMessage := fmt.Sprintf("Error in getSlug for %s: %s", shortUrl, err)
		c.HTML(http.StatusNotFound, "404.html", gin.H{"error": errorMessage})
		return
	}
	if _, err = s.IncreaseVisitForUrl(url, c.ClientIP()); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	c.HTML(http.StatusNotFound, "preview.tmpl", gin.H{"original": url.Original, "shortUrl": url.ShortURL()})
}

func dbMiddleware(s *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("store", s)
		c.Next()
	}
}