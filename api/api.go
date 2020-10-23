package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := NewRouter()
	router.ServeHTTP(w, r)
}

func NewRouter() http.Handler {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	for _, r := range cfg.Routes {
		for _, m := range r.Methods {
			switch m.Method {
			case http.MethodGet:
				router.GET(r.Path, func(c *gin.Context) {
					body, err := readFile(m.Response.File)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					//TODO m.Response.ContentType
					c.JSON(m.Response.StatusCode, body)
				})
			case http.MethodPost:
				router.POST(r.Path, func(c *gin.Context) {
					body, err := readFile(m.Response.File)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					//TODO m.Response.ContentType
					c.JSON(m.Response.StatusCode, body)
				})
			case http.MethodPut:
				router.PUT(r.Path, func(c *gin.Context) {
					body, err := readFile(m.Response.File)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					//TODO m.Response.ContentType
					c.JSON(m.Response.StatusCode, body)
				})
			case http.MethodPatch:
				router.PATCH(r.Path, func(c *gin.Context) {
					body, err := readFile(m.Response.File)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					//TODO m.Response.ContentType
					c.JSON(m.Response.StatusCode, body)
				})
			case http.MethodDelete:
				router.DELETE(r.Path, func(c *gin.Context) {
					body, err := readFile(m.Response.File)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					//TODO m.Response.ContentType
					c.JSON(m.Response.StatusCode, body)
				})
			}
		}
	}
	return router
}

func readFile(path string) ([]byte, error) {
	if len(path) > 0 {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("unable to read file %s: %w", path, err)
		}
		return b, nil
	}
	return nil, nil
}
