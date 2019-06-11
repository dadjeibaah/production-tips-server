package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dadjeibaah/production-tips-server/pkg/cache"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

var (
	query = flag.String("query", "ableton production tips", "Search term")
	port  = flag.String("port", "8080", "address to run server on")
)

func main() {
	flag.Parse()
	apiKey := os.Getenv("API_KEY")
	vimeoToken := os.Getenv("VIMEO_TOKEN")

	if *port == "" {
		log.Fatal("$PORT must be set")
	}

	if apiKey == "" {
		log.Fatal("$API_KEY must be set")
	}

	cch, err := cache.NewTipCache()
	if err != nil {
		log.Fatal(err)
	}

	cch = cch.WithVimeoSearcher(vimeoToken)

	router := gin.New()
	router.Use(gin.Logger())

	go cch.ClearSuggestionsOnTimeout()
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, cch.Search(*query))
	})

	router.Run(fmt.Sprintf(":%s", *port))
}
