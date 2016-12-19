package blimit_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/blimit"
	"github.com/insionng/macross/skipper"
	"testing"
)

func TestBodyLimit(t *testing.T) {
	e := macross.New()
	e.Use(blimit.BodyLimit("2M"))
	go e.Run(":6666")

	m := macross.New()
	m.Use(blimit.BodyLimitWithConfig(blimit.BodyLimitConfig{Skipper: skipper.DefaultSkipper, Limit: "4M"}))
	go m.Run(":7777")
}
