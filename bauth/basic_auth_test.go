package bauth_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/bauth"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	m := macross.New()
	m.Use(bauth.BasicAuth(func(username, password string) bool {
		if username == "inson" && password == "secret" {
			return true
		}
		return false
	}))

	go m.Run(":9999")
}
