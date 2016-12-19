package cors_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/cors"
	"github.com/insionng/macross/logger"
	"github.com/insionng/macross/recover"
	"testing"
)

func TestCORS(t *testing.T) {
	e := macross.New()
	e.Use(cors.CORSWithConfig(cors.CORSConfig{
		AllowOrigins: []string{"https://yougam.com", "https://github.com"},
		AllowHeaders: []string{macross.HeaderOrigin, macross.HeaderContentType, macross.HeaderAcceptEncoding},
	}))

	go e.Run(":8000")

}

var (
	users = []string{"Joe", "Veer", "Zion"}
)

func getUsers(c *macross.Context) error {
	return c.JSON(users)
}

func main() {
	e := macross.New()

	e.Use(logger.Logger())
	e.Use(recover.Recover())

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	// e.Use(cors.CORS())

	// CORS restricted
	// Allows requests from any `https://insionng.com` or `https://insionng.net` origin
	// wth GET, PUT, POST or DELETE method.
	e.Use(cors.CORSWithConfig(cors.CORSConfig{
		AllowOrigins: []string{"https://yougam.com", "https://github.com"},
		AllowMethods: []string{macross.GET, macross.PUT, macross.POST, macross.DELETE},
	}))

	e.Get("/api/users", getUsers)
	go e.Run(":9000")

}
