package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func watchlistHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Watchlist",
	})
}

func postHandler(c *gin.Context) {

	s := new(savereddit)
	s.Select()

	c.JSON(http.StatusOK, &s)
}

func Start() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	r.Use(static.Serve("/", static.LocalFile("./dashboard/build", true)))

	r.GET("/posts", postHandler)
	r.GET("/watchlist", watchlistHandler)

	if err := r.Run(":8080"); err != nil {
		Log.Error().Err(err).Send()
	}
}
