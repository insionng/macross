package pongor_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/pongor"
	"testing"
)

func TestRender(t *testing.T) {
	e := macross.New()
	e.SetRenderer(pongor.Renderor())
	e.Get("/", func() macross.Handler {
		return func(ctx *macross.Context) error {
			ctx.Set("title", "你好，世界")
			// render ./templates/index file.
			return ctx.Render("index")
		}
	}())
}
