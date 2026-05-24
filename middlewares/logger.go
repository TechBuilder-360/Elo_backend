package middlewares

import (
	"bytes"
	"time"

	"github.com/Toflex/directory_v2/pkg/log"
	"github.com/Toflex/directory_v2/pkg/util"
	"github.com/gin-gonic/gin"
)

type RequestLogger struct {
	RequestID string
	Status    int
	Latency   time.Duration
	Method    string
	Path      string
	Request   string
	Response  string
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)                  // capture
	return w.ResponseWriter.Write(b) // forward to client
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(log.SetLoggerInContext(c.Request.Context()))

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw

		c.Next()

		logger := log.LoggerInContext(c.Request.Context())
		logger.WithFields(map[string]interface{}{
			"request_id": c.GetString("request_id"),
			"status":     c.Writer.Status(),
			"latency":    time.Since(c.GetTime("start")),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"request":    util.ToJSON(c.Request.Body),
			"response":   util.ToJSON(blw.body.String()),
		}).Info("request completed")
	}
}
