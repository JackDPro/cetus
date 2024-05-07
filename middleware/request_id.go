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
		requestId := c.GetHeader("HTTP_X_REQUEST_ID")
		if requestId == "" {
			requestId = c.GetHeader("HTTP_REQUEST_ID")
		}
		if requestId == "" {
			requestId = uuid.New().String()
		}
		c.Set("request_id", requestId)
		c.Next()

		switch status := c.Writer.Status(); {
		case status < 400:
			_ = level.Info(provider.GetLogger()).Log("message", "success", "request_id", requestId)
		case 400 <= status && status < 500:
			_ = level.Warn(provider.GetLogger()).Log("path", c.Request.URL.String(), "params", c.Request.URL.Query(), "payload", c.Request.PostForm, "request_id", requestId)
		case 500 <= status:
			_ = level.Error(provider.GetLogger()).Log("path", c.Request.URL.String(), "params", c.Request.URL.Query(), "payload", c.Request.PostForm, "response", writer.body.String(), "request_id", requestId)
		}
	}
}
