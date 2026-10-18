package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "Login Success"})
		})
		auth.POST("/register", func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "Register Success"})
		})
	}

	return r
}
