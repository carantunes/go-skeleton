package middleware

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	logger "github.com/uniplaces/go-logger"
)

const successMessage = "success"

// Logger is our custom middleware for logging in a structured way
func Logger() (gin.HandlerFunc, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	localIP := GetLocalIP()

	handler := func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()
		contentType := c.Request.Header.Get("content-type")

		fields := map[string]interface{}{
			"hostname":      hostname,
			"local-ip":      localIP,
			"method":        method,
			"path":          path,
			"content-type":  contentType,
			"latency":       fmt.Sprintf("%13v", latency),
			"response-time": latency,
			"client-ip":     clientIP,
			"status":        statusCode,
			"user-agent":    userAgent,
			"errors":        handleErrorsField(c),
		}

		var lastError error

		if len(c.Errors) > 0 {
			lastError = c.Errors.Last()
		}

		logRequest(fields, lastError)
	}

	return handler, nil
}

func logRequest(fields map[string]interface{}, err error) {
	builder := logger.Builder().AddFields(fields)

	if err != nil {
		builder.Error(err)

		return
	}

	builder.Info(successMessage)
}

func handleErrorsField(ctx *gin.Context) []string {
	res := make([]string, len(ctx.Errors))
	for i, e := range ctx.Errors {
		res[i] = e.Error()
	}

	return res
}

// GetLocalIP returns the non loopback local IP of the host
// http://stackoverflow.com/a/31551220/916440
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
