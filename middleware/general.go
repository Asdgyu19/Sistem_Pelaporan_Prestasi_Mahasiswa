package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware untuk handle cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLogger middleware untuk logging requests
func RequestLogger() gin.HandlerFunc {
	return gin.Logger()
}

// Recovery middleware untuk handle panic
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
