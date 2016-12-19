package slash

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/skipper"
)

type (
	// TrailingSlashConfig defines the config for TrailingSlash middleware.
	TrailingSlashConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Status code to be used when redirecting the request.
		// Optional, but when provided the request is redirected using this code.
		RedirectCode int `json:"redirect_code"`
	}
)

var (
	// DefaultTrailingSlashConfig is the default TrailingSlash middleware config.
	DefaultTrailingSlashConfig = TrailingSlashConfig{
		Skipper: skipper.DefaultSkipper,
	}
)

// AddTrailingSlash returns a root level (before router) middleware which adds a
// trailing slash to the request `URL#Path`.
//
// Usage `Vodka#Pre(AddTrailingSlash())`
func AddTrailingSlash() macross.Handler {
	return AddTrailingSlashWithConfig(DefaultTrailingSlashConfig)
}

// AddTrailingSlashWithConfig returns a AddTrailingSlash middleware with config.
// See `AddTrailingSlash()`.
func AddTrailingSlashWithConfig(config TrailingSlashConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		//url := req.URL()
		path := string(c.Path())     //url.Path()
		qs := c.QueryArgs().String() // url.QueryString()
		if path != "/" && path[len(path)-1] != '/' {
			path += "/"
			uri := path
			if qs != "" {
				uri += "?" + qs
			}

			// Redirect
			if config.RedirectCode != 0 {
				return c.Redirect(uri, config.RedirectCode)
			}

			// Forward
			req.SetRequestURI(uri)
			req.URI().SetPath(path)
		}
		return c.Next()
	}
}

// RemoveTrailingSlash returns a root level (before router) middleware which removes
// a trailing slash from the request URI.
//
// Usage `Vodka#Pre(RemoveTrailingSlash())`
func RemoveTrailingSlash() macross.Handler {
	return RemoveTrailingSlashWithConfig(TrailingSlashConfig{})
}

// RemoveTrailingSlashWithConfig returns a RemoveTrailingSlash middleware with config.
// See `RemoveTrailingSlash()`.
func RemoveTrailingSlashWithConfig(config TrailingSlashConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		url := req.URI()
		path := string(url.Path())
		qs := c.QueryArgs().String() // url.QueryString()
		l := len(path) - 1
		if l >= 0 && path != "/" && path[l] == '/' {
			path = path[:l]
			uri := path
			if qs != "" {
				uri += "?" + qs
			}

			// Redirect
			if config.RedirectCode != 0 {
				c.Redirect(uri, config.RedirectCode)
				return nil
			}

			// Forward
			req.SetRequestURI(uri)
			req.URI().SetPath(path)
		}
		return c.Next()
	}
}
