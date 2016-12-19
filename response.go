package macross

import (
	"github.com/valyala/fasthttp"
)

// SetCookie implements `Context#SetCookie` function.
func (ctx *Context) SetCookie(c *Cookie) {
	cookie := new(fasthttp.Cookie)
	cookie.SetKey(c.Name())
	cookie.SetValue(c.Value())
	cookie.SetPath(c.Path())
	cookie.SetDomain(c.Domain())
	cookie.SetExpire(c.Expire())
	cookie.SetSecure(c.Secure())
	cookie.SetHTTPOnly(c.HTTPOnly())
	ctx.Response.Header.SetCookie(cookie)
}
