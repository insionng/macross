package compress

import (
	"compress/gzip"
	"github.com/insionng/macross"
	"github.com/insionng/macross/skipper"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type (
	// GzipConfig defines the config for Gzip middleware.
	GzipConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Gzip compression level.
		// Optional. Default value -1.
		Level int `json:"level"`
	}

	gzipResponseWriter struct {
		fasthttp.Response
		io.Writer
	}
)

var (
	// DefaultGzipConfig is the default Gzip middleware config.
	DefaultGzipConfig = GzipConfig{
		Skipper: skipper.DefaultSkipper,
		Level:   -1,
	}
)

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip() macross.Handler {
	return GzipWithConfig(DefaultGzipConfig)
}

// GzipWithConfig return Gzip middleware with config.
// See: `Gzip()`.
func GzipWithConfig(config GzipConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultGzipConfig.Skipper
	}
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}

	pool := gzipPool(config)
	scheme := "gzip"

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		res := c.Response
		res.Header.Add(macross.HeaderVary, macross.HeaderAcceptEncoding)
		if strings.Contains(string(c.Request.Header.Peek(macross.HeaderAcceptEncoding)), scheme) {
			rw := res.BodyWriter()
			gw := pool.Get().(*gzip.Writer)
			gw.Reset(rw)
			defer func() {
				if len(res.Body()) == 0 {
					// We have to reset response to it's pristine state when
					// nothing is written to body or error is returned.
					// See issue #424, #407.
					res.BodyWriteTo(rw)
					res.Header.Del(macross.HeaderContentEncoding)
					gw.Reset(ioutil.Discard)
				}
				gw.Close()
				pool.Put(gw)
			}()
			g := gzipResponseWriter{Response: res, Writer: gw}
			res.Header.Set(macross.HeaderContentEncoding, scheme)
			res.BodyWriteTo(g)
		}
		return c.Next()
	}
}

func (g gzipResponseWriter) Write(b []byte) (int, error) {
	if g.Header.Peek(macross.HeaderContentType) == nil {
		g.Header.Set(macross.HeaderContentType, http.DetectContentType(b))
	}
	return g.Writer.Write(b)
}

func gzipPool(config GzipConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(ioutil.Discard, config.Level)
			return w
		},
	}
}
