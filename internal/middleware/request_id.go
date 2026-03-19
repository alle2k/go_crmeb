package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := uuid.New()
		// Store in context as "request_id"
		c.Set("request_id", u.String())
		// Add to response header as "X-Request-ID"
		c.Header("X-Request-Id", u.String())
		c.Next()
	}
}
