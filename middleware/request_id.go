package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
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
		requestId := c.GetHeader("HTTP_X_REQUEST_ID")
		if requestId == "" {
			requestId = c.GetHeader("HTTP_REQUEST_ID")
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Set("request_id", requestId)
		c.Next()
	}
}
