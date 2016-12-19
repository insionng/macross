package logger_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	// Note: Just for the test coverage, not a real test.
	e := macross.New()
	e.Use(logger.LoggerWithConfig(logger.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	go e.Run(":9000")
}
