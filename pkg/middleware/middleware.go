package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-finance-dataset/internal/application/constants"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // capture response body
	return w.ResponseWriter.Write(b)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Respond to OPTIONS requests and stop further processing
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func ZapLoggerWithBody(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Clone the request body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // restore
		}

		// Wrap the response writer to capture response body
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Process request
		c.Next()

		end := time.Now()
		duration := end.Sub(start)
		formattedStart := start.Format("2006-01-02 15:04:05.000")
		formattedEnd := end.Format("2006-01-02 15:04:05.000")

		fields := []zap.Field{
			zap.String("start_time", formattedStart),
			zap.String("end_time", formattedEnd),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("latency", duration.String()),
			zap.ByteString("request_body", reqBody),
			zap.ByteString("response_body", blw.body.Bytes()),
		}

		// Add access_token if Authorization header is present
		if token := c.GetHeader(constants.Authorization); token != "" {
			fields = append(fields, zap.String("access_token", token))
		}

		// Add signature if Signature header is present
		if signature := c.GetHeader(constants.Signature); signature != "" {
			fields = append(fields, zap.String("signature", signature))
		}

		// Write structured log
		logger.Info("HTTP request", fields...)
	}
}
