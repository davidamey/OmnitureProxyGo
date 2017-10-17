package api

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/davidamey/omnitureproxy/archive"
	"github.com/gin-gonic/gin"
)

var rgxDate = regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}$")
var rgxVid = regexp.MustCompile("^\\d{38}$")

var fetcher archive.Reader = archive.NewReader("_archive")

func NewApi() http.Handler {
	r := gin.New()
	api := r.Group("/api")

	api.Any("/", func(c *gin.Context) {
		c.String(200, "/api/logs/")
	})

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
