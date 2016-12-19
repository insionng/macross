package gonder_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/gonder"
	"testing"
)

func TestRender(t *testing.T) {
	e := macross.New()
	e.SetRenderer(gonder.Renderor())
	e.Get("/", func() macross.Handler {
		return func(self *macross.Context) error {
			self.Set("title", "你好，世界")
			// render ./templates/index file.
			return self.Render("index")
		}
	}())
}
