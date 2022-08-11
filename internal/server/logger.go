package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger instances a Logger middleware for Gin.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		// clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		// Use debug level to keep production logs clean.
		zap.L().Debug(fmt.Sprintf("server: %s %s (%3d) [%v]",
			method,
			Log(path),
			statusCode,
			latency,
		))
	}
}

// Log sanitizes strings created from user input in response to the log4j debacle.
func Log(s string) string {
	if s == "" {
		return "''"
	} else if reject(s, 512) {
		return "?"
	}

	spaces := false

	// Remove non-printable and other potentially problematic characters.
	s = strings.Map(func(r rune) rune {
		if r < 32 {
			return -1
		}

		switch r {
		case ' ':
			spaces = true
			return r
		case '`', '"':
			return '\''
		case '\\', '$', '<', '>', '{', '}':
			return '?'
		default:
			return r
		}
	}, s)

	// Contains spaces?
	if spaces {
		return fmt.Sprintf("'%s'", s)
	}

	return s
}

// LogLower sanitizes strings created from user input and converts them to lowercase.
func LogLower(s string) string {
	return Log(strings.ToLower(s))
}
