package main

import (
	"asset-control/src/main/adapter/inbound/http/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	controller.Routes(r)

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
