package api

import (
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "X-Requested-With,X-Auth-Token,X-Humpy-Api-Key,Content-Type,Content-Length,Authorization")
		c.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}
