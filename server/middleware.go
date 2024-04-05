package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Logger instance a Logger middleware with config.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)

		if raw != "" {
			path = path + "?" + raw
		}

		logger := log.Info()
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			logger = log.Error().Any("error", c.Errors.Errors())
		}
		logger.
			Str("client_ip", c.ClientIP()).
			Str("path", path).
			Str("method", c.Request.Method).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Dur("latency", latency).
			Msg("request completed")

	}
}
