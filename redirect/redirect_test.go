package redirect_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/redirect"
	"testing"
)

func TestHTTPSRedirect(t *testing.T) {
	e := macross.New()
	e.Use(redirect.HTTPSRedirect())
	go e.Run(":6999")
}

func TestHTTPSWWWRedirect(t *testing.T) {
	e := macross.New()
	e.Use(redirect.HTTPSWWWRedirect())
	go e.Run(":7999")
}

func TestWWWRedirect(t *testing.T) {
	e := macross.New()
	e.Use(redirect.WWWRedirect())
	go e.Run(":8999")
}

func TestNonWWWRedirect(t *testing.T) {
	e := macross.New()
	e.Use(redirect.NonWWWRedirect())
	go e.Run(":9999")
}
