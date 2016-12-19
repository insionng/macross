package cors

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/skipper"
	"strconv"
	"strings"
)

type (
	// CORSConfig defines the config for CORS middleware.
	CORSConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// AllowOrigin defines a list of origins that may access the resource.
		// Optional. Default value []string{"*"}.
		AllowOrigins []string `json:"allow_origins"`

		// AllowMethods defines a list methods allowed when accessing the resource.
		// This is used in response to a preflight request.
		// Optional. Default value DefaultCORSConfig.AllowMethods.
		AllowMethods []string `json:"allow_methods"`

		// AllowHeaders defines a list of request headers that can be used when
		// making the actual request. This in response to a preflight request.
		// Optional. Default value []string{}.
		AllowHeaders []string `json:"allow_headers"`

		// AllowCredentials indicates whether or not the response to the request
		// can be exposed when the credentials flag is true. When used as part of
		// a response to a preflight request, this indicates whether or not the
		// actual request can be made using credentials.
		// Optional. Default value false.
		AllowCredentials bool `json:"allow_credentials"`

		// ExposeHeaders defines a whitelist headers that clients are allowed to
		// access.
		// Optional. Default value []string{}.
		ExposeHeaders []string `json:"expose_headers"`

		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached.
		// Optional. Default value 0.
		MaxAge int `json:"max_age"`
	}
)

var (
	// DefaultCORSConfig is the default CORS middleware config.
	DefaultCORSConfig = CORSConfig{
		Skipper:      skipper.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{macross.GET, macross.HEAD, macross.PUT, macross.PATCH, macross.POST, macross.DELETE},
	}
)

// CORS returns a Cross-Origin Resource Sharing (CORS) middleware.
// See: https://developer.mozilla.org/en/docs/Web/HTTP/Access_control_CORS
func CORS() macross.Handler {
	return CORSWithConfig(DefaultCORSConfig)
}

// CORSWithConfig returns a CORS middleware with config.
// See: `CORS()`.
func CORSWithConfig(config CORSConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultCORSConfig.Skipper
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		var origin = string(c.Request.Header.Peek(macross.HeaderOrigin))
		var allowOrigin string

		// Check allowed origins
		for _, o := range config.AllowOrigins {
			if o == "*" || o == origin {
				allowOrigin = o
				break
			}
		}

		// Simple request
		if string(c.Request.Header.Method()) != macross.OPTIONS {
			c.Response.Header.Add(macross.HeaderVary, macross.HeaderOrigin)
			c.Response.Header.Set(macross.HeaderAccessControlAllowOrigin, allowOrigin)
			if config.AllowCredentials {
				c.Response.Header.Set(macross.HeaderAccessControlAllowCredentials, "true")
			}
			if len(exposeHeaders) != 0 {
				c.Response.Header.Set(macross.HeaderAccessControlExposeHeaders, exposeHeaders)
			}
			return c.Next()
		}

		// Preflight request
		c.Response.Header.Add(macross.HeaderVary, macross.HeaderOrigin)
		c.Response.Header.Add(macross.HeaderVary, macross.HeaderAccessControlRequestMethod)
		c.Response.Header.Add(macross.HeaderVary, macross.HeaderAccessControlRequestHeaders)
		c.Response.Header.Set(macross.HeaderAccessControlAllowOrigin, allowOrigin)
		c.Response.Header.Set(macross.HeaderAccessControlAllowMethods, allowMethods)
		if config.AllowCredentials {
			c.Response.Header.Set(macross.HeaderAccessControlAllowCredentials, "true")
		}
		if len(allowHeaders) != 0 {
			c.Response.Header.Set(macross.HeaderAccessControlAllowHeaders, allowHeaders)
		} else {
			h := c.Request.Header.Peek(macross.HeaderAccessControlRequestHeaders)
			if h != nil {
				c.Response.Header.Set(macross.HeaderAccessControlAllowHeaders, string(h))
			}
		}
		if config.MaxAge > 0 {
			c.Response.Header.Set(macross.HeaderAccessControlMaxAge, maxAge)
		}
		return c.NoContent(macross.StatusNoContent)
	}
}
