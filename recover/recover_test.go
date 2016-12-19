package recover_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/recover"
	"testing"
)

func TestRecover(t *testing.T) {
	e := macross.New()
	e.Use(recover.RecoverWithConfig(recover.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	go e.Run(":8888")
}
