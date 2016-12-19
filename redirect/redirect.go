package redirect

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/skipper"
	"github.com/insionng/macross/slash"
	"net/http"
)

type (
	// RedirectConfig defines the config for Redirect middleware.
	RedirectConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Status code to be used when redirecting the request.
		// Optional. Default value http.StatusMovedPermanently.
		Code int `json:"code"`
	}
)

var (
	// DefaultRedirectConfig is the default Redirect middleware config.
	DefaultRedirectConfig = RedirectConfig{
		Skipper: skipper.DefaultSkipper,
		Code:    http.StatusMovedPermanently,
	}
)

// HTTPSRedirect redirects HTTP requests to HTTPS.
// For example, http://insionng.com will be redirect to https://insionng.com.
//
// Usage `Vodka#Pre(HTTPSRedirect())`
func HTTPSRedirect() macross.Handler {
	return HTTPSRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `HTTPSRedirect()`.
func HTTPSRedirectWithConfig(config RedirectConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = slash.DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		host := string(req.Host())
		uri := req.URI()
		if !c.RequestCtx.IsTLS() {
			c.Redirect("https://"+host+uri.String(), config.Code)
			return nil
		}
		return c.Next()
	}
}

// HTTPSWWWRedirect redirects HTTP requests to WWW HTTPS.
// For example, http://insionng.com will be redirect to https://www.insionng.com.
//
// Usage `Vodka#Pre(HTTPSWWWRedirect())`
func HTTPSWWWRedirect() macross.Handler {
	return HTTPSWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSWWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `HTTPSWWWRedirect()`.
func HTTPSWWWRedirectWithConfig(config RedirectConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = slash.DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		host := string(req.Host())
		uri := req.URI()
		if !c.RequestCtx.IsTLS() && host[:3] != "www" {
			c.Redirect("https://www."+host+uri.String(), http.StatusMovedPermanently)
			return nil
		}
		return c.Next()
	}
}

// WWWRedirect redirects non WWW requests to WWW.
// For example, http://insionng.com will be redirect to http://www.insionng.com.
//
// Usage `Vodka#Pre(WWWRedirect())`
func WWWRedirect() macross.Handler {
	return WWWRedirectWithConfig(DefaultRedirectConfig)
}

// WWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `WWWRedirect()`.
func WWWRedirectWithConfig(config RedirectConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = slash.DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		scheme := string(req.URI().Scheme())
		host := string(req.Host())
		if host[:3] != "www" {
			uri := req.URI()
			c.Redirect(scheme+"://www."+host+uri.String(), http.StatusMovedPermanently)
			return nil
		}
		return c.Next()
	}
}

// NonWWWRedirect redirects WWW requests to non WWW.
// For example, http://www.insionng.com will be redirect to http://insionng.com.
//
// Usage `Vodka#Pre(NonWWWRedirect())`
func NonWWWRedirect() macross.Handler {
	return NonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// NonWWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `NonWWWRedirect()`.
func NonWWWRedirectWithConfig(config RedirectConfig) macross.Handler {
	if config.Skipper == nil {
		config.Skipper = slash.DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		scheme := string(req.URI().Scheme())
		host := string(req.Host())
		if host[:3] == "www" {
			uri := req.URI()
			c.Redirect(scheme+"://"+host[4:]+uri.String(), http.StatusMovedPermanently)
			return nil
		}
		return c.Next()
	}

}
