package macross_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/access"
	"log"
	/*
		"github.com/insionng/macross/content"
		"github.com/insionng/macross/fault"
		"github.com/insionng/macross/file"
		"github.com/insionng/macross/slash"
	*/)

func MacrossTest() {
	m := macross.New()
	m.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
	)

	/*
		m.Use(
			// all these handlers are shared by every route
			access.Logger(log.Printf),
			slash.Remover(http.StatusMovedPermanently),
			fault.Recovery(log.Printf),
		)
	*/

	// serve RESTful APIs
	api := m.Group("/v1/api")
	/*
		api.Use(
			// these handlers are shared by the routes in the api group only
			content.TypeNegotiator(content.JSON, content.XML),
		)
	*/
	api.Get("/users", func(c *macross.Context) error {
		return c.Data("user list")
	})
	api.Post("/users", func(c *macross.Context) error {
		return c.Data("create a new user")
	})
	api.Put(`/users/<id:\d+>`, func(c *macross.Context) error {
		return c.Data("update user " + c.Param("id").String())
	})

	/*
		// serve index file
		m.Get("/", file.Content("ui/index.html"))
		// serve files under the "ui" subdirectory
		m.Get("/*", file.Server(file.PathMap{
			"/": "/ui/",
		}))
	*/

	m.Listen(":9000")

}
