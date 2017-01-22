// Package macross is a high productive and modular web framework in Golang.
package macross

import (
	ktx "context"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/valyala/fasthttp"
	"gopkg.in/ini.v1"
)

type (
	// Handler is the function for handling HTTP requests.
	Handler func(*Context) error

	// Macross manages routes and dispatches HTTP requests to the handlers of the matching routes.
	Macross struct {
		RouteGroup
		pool             sync.Pool
		routes           map[string]*Route
		stores           map[string]routeStore
		data             map[string]interface{} // data items managed by Key , Value
		maxParams        int
		binder           Binder
		sessioner        Sessioner
		localer          Localer
		notFound         []Handler
		notFoundHandlers []Handler
		renderer         Renderer
	}

	// routeStore stores route paths and the corresponding handlers.
	routeStore interface {
		Add(key string, data interface{}) int
		Get(key string, pvalues []string) (data interface{}, pnames []string)
		String() string
	}

	// Renderer is the interface that wraps the Render function.
	Renderer interface {
		Render(io.Writer, string, *Context) error
	}
)

// Export HTTP methods
const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

var (
	// Methods lists all supported HTTP methods by Macross.
	methods = [...]string{
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	}

	// Flash applies to current request.
	FlashNow bool

	// Configuration convention object.
	cfg *ini.File
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=utf-8"
)

// Headers
const (
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAllow                         = "Allow"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

// Status
// HTTP status codes as registered with IANA.
// See: http://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusPartialContent       = 206 // RFC 7233, 4.1
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
	StatusFound             = 302 // RFC 7231, 6.4.3
	StatusSeeOther          = 303 // RFC 7231, 6.4.4
	StatusNotModified       = 304 // RFC 7232, 4.1
	StatusUseProxy          = 305 // RFC 7231, 6.4.5
	_                       = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

var statusText = map[int]string{
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",
	StatusProcessing:         "Processing",

	StatusOK:                   "OK",
	StatusCreated:              "Created",
	StatusAccepted:             "Accepted",
	StatusNonAuthoritativeInfo: "Non-Authoritative Information",
	StatusNoContent:            "No Content",
	StatusResetContent:         "Reset Content",
	StatusPartialContent:       "Partial Content",
	StatusMultiStatus:          "Multi-Status",
	StatusAlreadyReported:      "Already Reported",
	StatusIMUsed:               "IM Used",

	StatusMultipleChoices:   "Multiple Choices",
	StatusMovedPermanently:  "Moved Permanently",
	StatusFound:             "Found",
	StatusSeeOther:          "See Other",
	StatusNotModified:       "Not Modified",
	StatusUseProxy:          "Use Proxy",
	StatusTemporaryRedirect: "Temporary Redirect",
	StatusPermanentRedirect: "Permanent Redirect",

	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusPaymentRequired:              "Payment Required",
	StatusForbidden:                    "Forbidden",
	StatusNotFound:                     "Not Found",
	StatusMethodNotAllowed:             "Method Not Allowed",
	StatusNotAcceptable:                "Not Acceptable",
	StatusProxyAuthRequired:            "Proxy Authentication Required",
	StatusRequestTimeout:               "Request Timeout",
	StatusConflict:                     "Conflict",
	StatusGone:                         "Gone",
	StatusLengthRequired:               "Length Required",
	StatusPreconditionFailed:           "Precondition Failed",
	StatusRequestEntityTooLarge:        "Request Entity Too Large",
	StatusRequestURITooLong:            "Request URI Too Long",
	StatusUnsupportedMediaType:         "Unsupported Media Type",
	StatusRequestedRangeNotSatisfiable: "Requested Range Not Satisfiable",
	StatusExpectationFailed:            "Expectation Failed",
	StatusTeapot:                       "I'm a teapot",
	StatusUnprocessableEntity:          "Unprocessable Entity",
	StatusLocked:                       "Locked",
	StatusFailedDependency:             "Failed Dependency",
	StatusUpgradeRequired:              "Upgrade Required",
	StatusPreconditionRequired:         "Precondition Required",
	StatusTooManyRequests:              "Too Many Requests",
	StatusRequestHeaderFieldsTooLarge:  "Request Header Fields Too Large",
	StatusUnavailableForLegalReasons:   "Unavailable For Legal Reasons",

	StatusInternalServerError:           "Internal Server Error",
	StatusNotImplemented:                "Not Implemented",
	StatusBadGateway:                    "Bad Gateway",
	StatusServiceUnavailable:            "Service Unavailable",
	StatusGatewayTimeout:                "Gateway Timeout",
	StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
	StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
	StatusInsufficientStorage:           "Insufficient Storage",
	StatusLoopDetected:                  "Loop Detected",
	StatusNotExtended:                   "Not Extended",
	StatusNetworkAuthenticationRequired: "Network Authentication Required",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}

// New creates a new Macross object.
func New() *Macross {
	m := &Macross{
		routes: make(map[string]*Route),
		stores: make(map[string]routeStore),
	}
	m.RouteGroup = *newRouteGroup("", m, make([]Handler, 0))
	m.NotFound(MethodNotAllowedHandler, NotFoundHandler)
	m.SetBinder(&binder{})
	m.pool.New = func() interface{} {
		return &Context{
			ktx:     ktx.Background(),
			pvalues: make([]string, m.maxParams),
			macross: m,
		}
	}
	return m
}

// New creates a new Classic Macross object.
func Classic() *Macross {
	m := &Macross{
		routes: make(map[string]*Route),
		stores: make(map[string]routeStore),
	}
	m.RouteGroup = *newRouteGroup("", m, make([]Handler, 0))
	m.NotFound(MethodNotAllowedHandler, NotFoundHandler)
	m.SetBinder(&binder{})
	m.SetSessioner(&sessioner{})
	m.SetLocaler(&localer{})
	m.pool.New = func() interface{} {
		return &Context{
			ktx:     ktx.Background(),
			Session: m.sessioner,
			Localer: m.localer,
			pvalues: make([]string, m.maxParams),
			macross: m,
		}
	}
	return m
}

// AcquireContext returns an empty `Context` instance from the pool.
// You must return the context by calling `ReleaseContext()`.
func (m *Macross) AcquireContext() *Context {
	if ctx, okay := m.pool.Get().(*Context); okay {
		return ctx
	} else {
		panic("Not Standard Macross Context")
		return nil
	}
}

// ReleaseContext returns the `Context` instance back to the pool.
// You must call it after `AcquireContext()`.
func (m *Macross) ReleaseContext(c *Context) {
	c.Response.Header.SetServer("Macross")
	m.pool.Put(c)
}

// ServeHTTP handles the HTTP request.
func (m *Macross) ServeHTTP(ctx *fasthttp.RequestCtx) {
	c := m.AcquireContext()
	c.Reset(ctx)
	c.handlers, c.pnames = m.find(string(ctx.Method()), string(ctx.Path()), c.pvalues)
	if err := c.Next(); err != nil {
		m.HandleError(c, err)
	}
	m.ReleaseContext(c)
}

// Route returns the named route.
// Nil is returned if the named route cannot be found.
func (r *Macross) Route(name string) *Route {
	return r.routes[name]
}

// Use appends the specified handlers to the macross and shares them with all routes.
func (r *Macross) Use(handlers ...Handler) {
	r.RouteGroup.Use(handlers...)
	r.notFoundHandlers = combineHandlers(r.handlers, r.notFound)
}

// NotFound specifies the handlers that should be invoked when the macross cannot find any route matching a request.
// Note that the handlers registered via Use will be invoked first in this case.
func (r *Macross) NotFound(handlers ...Handler) {
	r.notFound = handlers
	r.notFoundHandlers = combineHandlers(r.handlers, r.notFound)
}

// HandleError is the error handler for handling any unhandled errors.
func (m *Macross) HandleError(c *Context, err interface{}) {
	status := StatusInternalServerError
	msg := StatusText(status)
	if httpError, okay := err.(*HTTPError); okay {
		status = httpError.Status
		msg = httpError.Message
	} else if iError, okay := err.(error); okay {
		msg = iError.Error()
	}

	if string(c.Request.Header.Method()) == HEAD {
		c.NoContent(status)
	} else {
		c.String(msg, status)
	}
}

func (r *Macross) add(method, path string, handlers []Handler) {
	store := r.stores[method]
	if store == nil {
		store = newStore()
		r.stores[method] = store
	}
	if n := store.Add(path, handlers); n > r.maxParams {
		r.maxParams = n
	}
}

func (r *Macross) find(method, path string, pvalues []string) (handlers []Handler, pnames []string) {
	var hh interface{}
	if store := r.stores[method]; store != nil {
		hh, pnames = store.Get(path, pvalues)
	}
	if hh != nil {
		return hh.([]Handler), pnames
	}
	return r.notFoundHandlers, pnames
}

func (r *Macross) findAllowedMethods(path string) map[string]bool {
	methods := make(map[string]bool)
	pvalues := make([]string, r.maxParams)
	for m, store := range r.stores {
		if handlers, _ := store.Get(path, pvalues); handlers != nil {
			methods[m] = true
		}
	}
	return methods
}

func (m *Macross) Pull(key string) interface{} {
	return m.data[key]
}

func (m *Macross) Push(key string, value interface{}) {
	if m.data == nil {
		m.data = make(map[string]interface{})
	}
	m.data[key] = value
}

// SetBinder registers a custom binder. It's invoked by `Context#Bind()`.
func (m *Macross) SetBinder(b Binder) {
	m.binder = b
}

// Binder returns the binder instance.
func (m *Macross) Binder() Binder {
	return m.binder
}

// SetSessioner registers a custom sessioner. It's invoked by `Context#Bind()`.
func (m *Macross) SetSessioner(s Sessioner) {
	m.sessioner = s
}

// Sessioner returns the sessioner instance.
func (m *Macross) Sessioner() Sessioner {
	return m.sessioner
}

// SetLocale registers a custom localer.
func (m *Macross) SetLocaler(l Localer) {
	m.localer = l
}

// Localer returns the locale instance.
func (m *Macross) Localer() Localer {
	return m.localer
}

// SetRenderer registers an HTML template renderer. It's invoked by `Context#Render()`.
func (m *Macross) SetRenderer(r Renderer) {
	m.renderer = r
}

// Static registers a new route with path prefix to serve static files from the
// provided root directory.
func (m *Macross) Static(prefix, root string) {
	if prefix == "/" {
		prefix = prefix + "*"
	} else if len(prefix) > 1 {
		if prefix[:1] != "/" {
			prefix = prefix + "/*"
		} else {
			prefix = prefix + "*"
		}
	}
	m.Get(prefix, func(c *Context) error {
		return c.ServeFile(path.Join(root, c.Parameter(0)))
	})
}

// File registers a new route with path to serve a static file.
func (m *Macross) File(path, file string) {
	m.Get(path, func(c *Context) error {
		return c.ServeFile(file)
	})
}

// NotFoundHandler returns a 404 HTTP error indicating a request has no matching route.
func NotFoundHandler(*Context) error {
	return NewHTTPError(StatusNotFound)
}

// MethodNotAllowedHandler handles the situation when a request has matching route without matching HTTP method.
// In this case, the handler will respond with an Allow HTTP header listing the allowed HTTP methods.
// Otherwise, the handler will do nothing and let the next handler (usually a NotFoundHandler) to handle the problem.
func MethodNotAllowedHandler(c *Context) error {
	methods := c.Macross().findAllowedMethods(string(c.Path()))
	if len(methods) == 0 {
		return nil
	}
	methods["OPTIONS"] = true
	ms := make([]string, len(methods))
	i := 0
	for method := range methods {
		ms[i] = method
		i++
	}
	sort.Strings(ms)
	c.Response.Header.Set("Allow", strings.Join(ms, ", "))
	if string(c.Method()) != "OPTIONS" {
		c.Response.SetStatusCode(StatusMethodNotAllowed)
	}
	return c.Abort()
}

func GetAddress(args ...interface{}) string {

	var host string
	var port int

	if len(args) == 1 {
		switch arg := args[0].(type) {
		case string:
			addrs := strings.Split(args[0].(string), ":")
			if len(addrs) == 1 {
				host = addrs[0]
			} else if len(addrs) >= 2 {
				host = addrs[0]
				_port, _ := strconv.ParseInt(addrs[1], 10, 0)
				port = int(_port)
			}
		case int:
			port = arg
		}
	} else if len(args) >= 2 {
		if arg, ok := args[0].(string); ok {
			host = arg
		}
		if arg, ok := args[1].(int); ok {
			port = arg
		}
	}

	if host_ := os.Getenv("HOST"); len(host_) != 0 {
		host = host_
	} else if len(host) == 0 {
		host = "0.0.0.0"
	}

	if port_, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 32); port_ != 0 {
		port = int(port_)
	} else if port == 0 {
		port = 8000
	}

	addr := host + ":" + strconv.FormatInt(int64(port), 10)
	return addr

}

// SetConfig sets data sources for configuration.
func SetConfig(source interface{}, others ...interface{}) (_ *ini.File, err error) {
	cfg, err = ini.Load(source, others...)
	return Config(), err
}

// Config returns configuration convention object.
// It returns an empty object if there is no one available.
func Config() *ini.File {
	if cfg == nil {
		return ini.Empty()
	}
	return cfg
}
