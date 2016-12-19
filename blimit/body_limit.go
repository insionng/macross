package blimit

import (
	"bytes"
	"fmt"
	"github.com/insionng/macross"
	lbytes "github.com/insionng/macross/libraries/gommon/bytes"
	"github.com/insionng/macross/skipper"
	"io"
	"sync"
)

type (
	// BodyLimitConfig defines the config for BodyLimit middleware.
	BodyLimitConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Maximum allowed size for a request body, it can be specified
		// as `4x` or `4xB`, where x is one of the multiple from K, M, G, T or P.
		Limit string `json:"limit"`
		limit int64
	}

	limitedReader struct {
		BodyLimitConfig
		reader  io.Reader
		read    int64
		context *macross.Context
	}
)

var (
	// DefaultBodyLimitConfig is the default Gzip middleware config.
	DefaultBodyLimitConfig = BodyLimitConfig{
		Skipper: skipper.DefaultSkipper,
	}
)

// BodyLimit returns a BodyLimit middleware.
//
// BodyLimit middleware sets the maximum allowed size for a request body, if the
// size exceeds the configured limit, it sends "413 - Request Entity Too Large"
// response. The BodyLimit is determined based on both `Content-Length` request
// header and actual content read, which makes it super secure.
// Limit can be specified as `4x` or `4xB`, where x is one of the multiple from K, M,
// G, T or P.
func BodyLimit(limit string) macross.Handler {
	c := DefaultBodyLimitConfig
	c.Limit = limit
	return BodyLimitWithConfig(c)
}

// BodyLimitWithConfig returns a BodyLimit middleware with config.
// See: `BodyLimit()`.
func BodyLimitWithConfig(config BodyLimitConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultBodyLimitConfig.Skipper
	}

	limit, err := lbytes.Parse(config.Limit)
	if err != nil {
		panic(fmt.Errorf("invalid body-limit=%s", config.Limit))
	}
	config.limit = limit
	pool := limitedReaderPool(config)

	return func(c *macross.Context) error {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request

		// Based on content length
		if int64(req.Header.ContentLength()) > config.limit {
			return macross.ErrStatusRequestEntityTooLarge
		}

		// Based on content read
		r := pool.Get().(*limitedReader)
		r.Reset(bytes.NewBuffer(req.Body()), c)
		defer pool.Put(r)
		req.SetBodyStream(r, 0)

		return c.Next()
	}

}

func (r *limitedReader) Read(b []byte) (n int, err error) {
	n, err = r.reader.Read(b)
	r.read += int64(n)
	if r.read > r.limit {
		return n, macross.ErrStatusRequestEntityTooLarge
	}
	return
}

func (r *limitedReader) Reset(reader io.Reader, context *macross.Context) {
	r.reader = reader
	r.context = context
}

func limitedReaderPool(c BodyLimitConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			return &limitedReader{BodyLimitConfig: c}
		},
	}
}
