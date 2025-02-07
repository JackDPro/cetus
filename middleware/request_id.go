package middleware

import (
	"bytes"
	"github.com/JackDPro/cetus/provider"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		writer := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer
		_ = level.Info(provider.GetLogger()).Log("id", "jack", "header", c.Request.Header)
		requestId := c.GetHeader("HTTP_X_REQUEST_ID")
		if requestId == "" {
			requestId = c.GetHeader("HTTP_REQUEST_ID")
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Set("request_id", requestId)
		c.Writer.Header().Add("Request-Id", requestId)
		c.Next()
	}
}
