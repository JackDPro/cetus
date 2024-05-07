package middleware

import (
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
)

func LimitRate(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.AbortWithStatusJSON(httpError.StatusCode, &map[string]string{
				"message": httpError.Message,
			})
		} else {
			c.Next()
		}
	}
}
