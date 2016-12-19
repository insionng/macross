package main

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/fempla"
	"github.com/insionng/macross/logger"
	"github.com/insionng/macross/recover"
	"github.com/insionng/macross/static"
)

func main() {

	v := macross.New()
	v.Use(logger.Logger())
	v.Use(recover.Recover())
	v.SetRenderer(fempla.Renderor())
	v.Use(static.Static("static"))
	v.Get("/", func(self *macross.Context) error {
		data := make(map[string]interface{})
		data["oh"] = "no"
		data["name"] = "Insion Ng"
		self.Set("title", "你好，世界")
		self.SetStore(data)
		self.Set("oh", "yes")
		return self.Render("index")
	})

	v.Listen(9000)

}
