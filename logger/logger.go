package logger

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/insionng/macross"
	"github.com/insionng/macross/libraries/gommon/color"
	"github.com/insionng/macross/skipper"
	isatty "github.com/mattn/go-isatty"
	"github.com/valyala/fasttemplate"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper skipper.Skipper

		// Log format which can be constructed using the following tags:
		//
		// - time_rfc3339
		// - id (Request ID - Not implemented)
		// - remote_ip
		// - uri
		// - host
		// - method
		// - path
		// - referer
		// - user_agent
		// - status
		// - latency (In microseconds)
		// - latency_human (Human readable)
		// - bytes_in (Bytes received)
		// - bytes_out (Bytes sent)
		//
		// Example "${remote_ip} ${status}"
		//
		// Optional. Default value DefaultLoggerConfig.Format.
		Format string `json:"format"`

		// Output is a writer where logs are written.
		// Optional. Default value os.Stdout.
		Output io.Writer

		template   *fasttemplate.Template
		color      *color.Color
		bufferPool sync.Pool
	}
)

var (
	// DefaultLoggerConfig is the default Logger middleware config.
	DefaultLoggerConfig = LoggerConfig{
		Skipper: skipper.DefaultSkipper,
		Format: `{"time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
			`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
			`"bytes_out":${bytes_out}}` + "\n",
		Output: os.Stdout,
		color:  color.New(),
	}
)

// Logger returns a middleware that logs HTTP requests.
func Logger() macross.Handler {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) macross.Handler {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.Format == "" {
		config.Format = DefaultLoggerConfig.Format
	}
	if config.Output == nil {
		config.Output = DefaultLoggerConfig.Output
	}

	config.template = fasttemplate.New(config.Format, "${", "}")
	config.color = color.New()
	if w, ok := config.Output.(*os.File); !ok || !isatty.IsTerminal(w.Fd()) {
		config.color.Disable()
	}
	config.bufferPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 256))
		},
	}

	return func(c *macross.Context) (err error) {
		if config.Skipper(c) {
			return c.Next()
		}

		req := c.Request
		start := time.Now()
		if err = c.Next(); err != nil {
			if herr, okay := err.(*macross.HTTPError); okay {
				c.Error(herr.Message, herr.Status)
			}
		}
		stop := time.Now()
		buf := config.bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer config.bufferPool.Put(buf)

		_, err = config.template.ExecuteFunc(buf, func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "time_rfc3339":
				return w.Write([]byte(time.Now().Format(time.RFC3339)))
			case "remote_ip":
				ra := c.RealIP()
				return w.Write([]byte(ra))
			case "host":
				return w.Write(req.Host())
			case "uri":
				return w.Write(req.URI().FullURI())
			case "method":
				return w.Write(req.Header.Method())
			case "path":
				p := string(c.Path())
				if p == "" {
					p = "/"
				}
				return w.Write([]byte(p))
			case "referer":
				return w.Write(req.Header.Referer())
			case "user_agent":
				return w.Write(req.Header.UserAgent())
			case "status":
				n := c.Response.StatusCode()
				s := config.color.Green(n)
				switch {
				case n >= 500:
					s = config.color.Red(n)
				case n >= 400:
					s = config.color.Yellow(n)
				case n >= 300:
					s = config.color.Cyan(n)
				}
				return w.Write([]byte(s))
			case "latency":
				l := stop.Sub(start).Nanoseconds() / 1000
				return w.Write([]byte(strconv.FormatInt(l, 10)))
			case "latency_human":
				return w.Write([]byte(stop.Sub(start).String()))
			case "bytes_in":
				b := string(req.Header.Peek(macross.HeaderContentLength))
				if b == "" {
					b = "0"
				}
				return w.Write([]byte(b))
			case "bytes_out":
				res := c.Response
				size := int64(len(res.Body()))
				return w.Write([]byte(strconv.FormatInt(size, 10)))
			}
			return 0, nil
		})
		if err == nil {
			config.Output.Write(buf.Bytes())
		}
		return
	}
}
