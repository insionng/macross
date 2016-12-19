package slash_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/slash"
	"testing"
)

func TestTrailingSlash(t *testing.T) {
	e := macross.New()
	e.Use(slash.AddTrailingSlash())
	go e.Run(":8888")

	m := macross.New()
	m.Use(slash.RemoveTrailingSlash())
	go m.Run(":9999")
}
