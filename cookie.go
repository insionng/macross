package macross

import (
	"github.com/valyala/fasthttp"
	"time"
)

type (
	// Cookie implements `Cookie`.
	Cookie struct {
		fasthttp.Cookie
	}
)

// Name returns the cookie name.
func (c *Cookie) Name() string {
	return string(c.Cookie.Key())
}

// SetName sets cookie name.
func (c *Cookie) SetName(name string) {
	c.Cookie.SetKey(name)
}

// Value returns the cookie value.
func (c *Cookie) Value() string {
	return string(c.Cookie.Value())
}

// SetValue sets the cookie value.
func (c *Cookie) SetValue(value string) {
	c.Cookie.SetValue(value)
}

// Path returns the cookie path.
func (c *Cookie) Path() string {
	return string(c.Cookie.Path())
}

// SetPath sets the cookie path.
func (c *Cookie) SetPath(path string) {
	c.Cookie.SetPath(path)
}

// Domain returns the cookie domain.
func (c *Cookie) Domain() string {
	return string(c.Cookie.Domain())
}

// SetDomain sets the cookie domain.
func (c *Cookie) SetDomain(domain string) {
	c.Cookie.SetDomain(domain)
}

// Expires returns the cookie expiry time.
func (c *Cookie) Expire() time.Time {
	return c.Cookie.Expire()
}

// SetExpire sets the cookie expiry time.
func (c *Cookie) SetExpire(expire time.Time) {
	c.Cookie.SetExpire(expire)
}

// Secure indicates if cookie is Secure.
func (c *Cookie) Secure() bool {
	return c.Cookie.Secure()
}

// SetSecure sets the cookie as Secure.
func (c *Cookie) SetSecure(secure bool) {
	c.Cookie.SetSecure(secure)
}

// HTTPOnly indicates if cookie is HTTPOnly.
func (c *Cookie) HTTPOnly() bool {
	return c.Cookie.HTTPOnly()
}

// SetHTTPOnly sets the cookie as HTTPOnly.
func (c *Cookie) SetHTTPOnly(httpOnly bool) {
	c.Cookie.SetHTTPOnly(httpOnly)
}
