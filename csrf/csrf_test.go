package csrf_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/csrf"
	"testing"
)

func TestCSRF(t *testing.T) {
	e := macross.New()
	e.Use(csrf.CSRFWithConfig(csrf.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
	}))
	go e.Listen(9000)
}
