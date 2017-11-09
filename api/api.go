package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/davidamey/omnitureproxy/archive"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/gin-gonic/gin"
)

var rgxDate = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}$")
var rgxVid = regexp.MustCompile("^\\d{3,38}$") //todo: remove 3, once testing is done

var fetcher archive.Reader = archive.NewReader("_archive")

func NewApi() http.Handler {
	r := gin.New()
	r.RedirectTrailingSlash = false
	api := r.Group("/api")

	api.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-cache, must-revalidate")

		c.Next()
	})

	api.Any("/", func(c *gin.Context) {
		c.String(200, strings.Join([]string{
			"/api/logs",
			"/api/logs/<date>",
			"/api/logs/<date>/<visitor ID>",
		}, "\n"))
	})

	api.Any("/stats", gin.WrapF(stats_api.Handler))

	api.GET("/logs", func(c *gin.Context) {
		c.JSON(200, fetcher.GetDates())
	})

	api.GET("/logs/:date", func(c *gin.Context) {
		d := c.Param("date")
		if !rgxDate.MatchString(d) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid date"))
			return
		}
		c.JSON(200, fetcher.GetVisitorsForDate(d))
	})

	api.GET("/logs/:date/:vid", func(c *gin.Context) {
		d, v := c.Param("date"), c.Param("vid")
		if !rgxDate.MatchString(d) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid date"))
			return
		}
		if !rgxVid.MatchString(v) {
			c.AbortWithError(http.StatusBadRequest, fmt.Errorf("invalid vid"))
			return
		}
		c.Data(200, "application/json; charset=utf-8", fetcher.GetArchive(d, v))
	})

	return r
}
