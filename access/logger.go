// Package access provides an access logging handler for the macross package.
package access

import (
	"fmt"
	"github.com/insionng/macross"
	"time"
)

// LogFunc logs a message using the given format and optional arguments.
// The usage of format and arguments is similar to that for fmt.Printf().
// LogFunc should be thread safe.
type LogFunc func(format string, a ...interface{})

// Logger returns a handler that logs a message for every request.
// The access log messages contain information including client IPs, time used to serve each request, request line,
// response status and size.
//
//     import (
//         "log"
//         "macross"
//         "macross/access"
//     )
//
//     r := macross.New()
//     r.Use(access.Logger(log.Printf))
func Logger(log LogFunc) macross.Handler {
	return func(c *macross.Context) error {
		startTime := time.Now()

		err := c.Next()

		clientIP := c.RemoteIP()
		elapsed := float64(time.Now().Sub(startTime).Nanoseconds()) / 1e6
		requestLine := fmt.Sprintf("%s %s", c.Method(), c.URI().String())
		log(`[%s] [%.3fms] %s %d`, clientIP, elapsed, requestLine, c.Response.StatusCode())
		fmt.Println("\n")
		return err
	}
}

/*
func getClientIP(req *http.Request) string {
	ip := req.Header.Get("X-Real-IP")
	if ip == "" {
		ip = req.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = req.RemoteAddr
		}
	}
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
*/
