package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(r *gin.Engine) {
	api := r.Group("/api")

	// Router to confirm the condition of server
	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	api.POST("/employee", //Use case employee)
}
