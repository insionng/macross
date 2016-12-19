package moverride_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/moverride"
	"testing"
)

func TestMethodOverride(t *testing.T) {
	e := macross.New()
	e.Use(moverride.MethodOverrideWithConfig(moverride.MethodOverrideConfig{
		Getter: moverride.MethodFromForm("_method"),
	}))
	go e.Run(":9000")
}
