package moverride

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/skipper"
)

type (
	// MethodOverrideConfig defines the config for MethodOverride middleware.
	MethodOverrideConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Getter is a function that gets overridden method from the request.
		// Optional. Default values MethodFromHeader(macross.HeaderXHTTPMethodOverride).
		Getter MethodOverrideGetter
	}

	// MethodOverrideGetter is a function that gets overridden method from the request
	MethodOverrideGetter func(*macross.Context) string
)

var (
	// DefaultMethodOverrideConfig is the default MethodOverride middleware config.
	DefaultMethodOverrideConfig = MethodOverrideConfig{
		Skipper: skipper.DefaultSkipper,
		Getter:  MethodFromHeader(macross.HeaderXHTTPMethodOverride),
	}
)

// MethodOverride returns a MethodOverride middleware.
// MethodOverride  middleware checks for the overridden method from the request and
// uses it instead of the original method.
//
// For security reasons, only `POST` method can be overridden.
func MethodOverride() macross.Handler {
	return MethodOverrideWithConfig(DefaultMethodOverrideConfig)
}

// MethodOverrideWithConfig returns a MethodOverride middleware with config.
// See: `MethodOverride()`.
func MethodOverrideWithConfig(config MethodOverrideConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultMethodOverrideConfig.Skipper
	}
	if config.Getter == nil {
		config.Getter = DefaultMethodOverrideConfig.Getter
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		if string(req.Header.Method()) == macross.POST {
			m := config.Getter(c)
			if m != "" {
				req.Header.SetMethod(m)
			}
		}
		return c.Next()
	}
}

// MethodFromHeader is a `MethodOverrideGetter` that gets overridden method from
// the request header.
func MethodFromHeader(header string) MethodOverrideGetter {
	return func(c *macross.Context) string {
		return string(c.Request.Header.Peek(header))
	}
}

// MethodFromForm is a `MethodOverrideGetter` that gets overridden method from the
// form parameter.
func MethodFromForm(param string) MethodOverrideGetter {
	return func(c *macross.Context) string {
		return string(c.FormValue(param))
	}
}

// MethodFromQuery is a `MethodOverrideGetter` that gets overridden method from
// the query parameter.
func MethodFromQuery(param string) MethodOverrideGetter {
	return func(c *macross.Context) string {
		return string(c.QueryArgs().Peek(param))
	}
}
