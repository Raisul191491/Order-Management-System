package main

import (
	"oms/container"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.New()
	e.Use(gin.Recovery())
	container.Serve(e)
}
