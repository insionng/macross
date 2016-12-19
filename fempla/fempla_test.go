package fempla_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/fempla"
	"testing"
)

func TestRender(t *testing.T) {
	m := macross.New()
	m.SetRenderer(fempla.Renderor())
	m.Get("/", func() macross.Handler {
		return func(self *macross.Context) error {
			self.Set("title", "你好，世界")

			// render ./templates/index.html file.
			return self.Render("index")
		}
	}())
}
