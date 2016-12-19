package bauth

import (
	"encoding/base64"
	"github.com/insionng/macross"
	"github.com/insionng/macross/skipper"
)

type (
	// BasicAuthConfig defines the config for BasicAuth middleware.
	BasicAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Validator is a function to validate BasicAuth credentials.
		// Required.
		Validator BasicAuthValidator
	}

	// BasicAuthValidator defines a function to validate BasicAuth credentials.
	BasicAuthValidator func(string, string) bool
)

const (
	basic = "Basic"
)

var (
	// DefaultBasicAuthConfig is the default BasicAuth middleware config.
	DefaultBasicAuthConfig = BasicAuthConfig{
		Skipper: skipper.DefaultSkipper,
	}
)

// BasicAuth returns an BasicAuth middleware.
//
// For valid credentials it calls the next handler.
// For invalid credentials, it sends "401 - Unauthorized" response.
// For empty or invalid `Authorization` header, it sends "400 - Bad Request" response.
func BasicAuth(fn BasicAuthValidator) macross.Handler {
	c := DefaultBasicAuthConfig
	c.Validator = fn
	return BasicAuthWithConfig(c)
}

// BasicAuthWithConfig returns an BasicAuth middleware with config.
// See `BasicAuth()`.
func BasicAuthWithConfig(config BasicAuthConfig) macross.Handler {
	// Defaults
	if config.Validator == nil {
		panic("basic-auth middleware requires validator function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBasicAuthConfig.Skipper
	}

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		auth := string(c.Request.Header.Peek(macross.HeaderAuthorization))
		l := len(basic)

		if len(auth) > l+1 && auth[:l] == basic {
			b, err := base64.StdEncoding.DecodeString(auth[l+1:])
			if err != nil {
				return err
			}
			cred := string(b)
			for i := 0; i < len(cred); i++ {
				if cred[i] == ':' {
					// Verify credentials
					if config.Validator(cred[:i], cred[i+1:]) {
						return c.Next()
					}
				}
			}
		}
		// Need to return `401` for browsers to pop-up login box.
		c.Response.Header.Set(macross.HeaderWWWAuthenticate, basic+" realm=Restricted")
		return macross.ErrUnauthorized
	}

}
