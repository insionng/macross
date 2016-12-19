package main

import (
	"fmt"
	"github.com/insionng/macross"
	"github.com/insionng/macross/access"
	"log"
	/*
		"github.com/insionng/macross/content"
		"github.com/insionng/macross/fault"
		"github.com/insionng/macross/file"
		"github.com/insionng/macross/slash"
	*/)

type User struct {
	Name  string
	Email string
}

func Auther() macross.Handler {
	return func(self *macross.Context) error {
		self.Set("version", "1.0.0")
		u := User{Name: "Insion Ng", Email: "insion@live.com"}
		self.Set("user", u)
		return self.Next()
	}
}

func Checker() macross.Handler {
	return func(self *macross.Context) error {
		fmt.Println("Checker() macross.Handler>", self.Get("user"))
		if u := self.Get("user"); u != nil {
			fmt.Println("Ck:", u)
			return self.Next()
		} else {
			self.Redirect("/signin/", macross.StatusFound)
			return nil
		}

	}
}

func main() {
	m := macross.New()
	m.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		Auther(),
	)

	m.Any("/signin/", func(self *macross.Context) error {
		return self.Data("signin")
	})

	g := m.Group("", Checker())
	g.Get("/", func(self *macross.Context) error {
		fmt.Println(self.Get("version"))
		fmt.Println(self.Get("user"))
		return self.Data("Macross")
	})

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
	api.Get("/users", func(self *macross.Context) error {
		fmt.Println(self.Get("version"))
		fmt.Println(self.Get("user"))
		return self.Data("user list")
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

	m.Listen(8080)
}
