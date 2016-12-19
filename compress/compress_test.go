package compress_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/compress"
	"testing"
)

func TestGzip(t *testing.T) {
	e := macross.New()
	e.Use(compress.GzipWithConfig(compress.GzipConfig{
		Level: 5,
	}))
	go e.Run(":9000")
}
