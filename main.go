package main

import (
	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	content := gin.H{"Hello": "World"}
	c.JSON(200, content)
}

func main() {
	router := gin.Default()
	router.GET("/", index)
	router.Run()
}
