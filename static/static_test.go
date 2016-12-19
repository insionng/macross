package static_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/static"
	"testing"
)

func TestStatic(t *testing.T) {
	m := macross.New()
	m.Use(static.Static("public"))
	go m.Run(":8888")

	n := macross.New()
	n.Use(static.StaticWithConfig(static.StaticConfig{
		Root:   "public",
		Browse: true,
	}))
	go n.Run(":9999")
}
