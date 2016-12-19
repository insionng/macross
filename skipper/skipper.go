package skipper

import "github.com/insionng/macross"

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c *macross.Context) bool
)

// defaultSkipper returns false which processes the middleware.
func DefaultSkipper(c *macross.Context) bool {
	return false
}
