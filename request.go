package macross

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"

	"github.com/valyala/fasthttp"
)

// IsTLS implements `Context#TLS` function.
func (c *Context) IsTLS() bool {
	return c.RequestCtx.IsTLS()
}

// Scheme implements `Context#Scheme` function.
func (c *Context) Scheme() string {
	return string(c.RequestCtx.URI().Scheme())
}

// Host implements `Context#Host` function.
func (c *Context) Host() string {
	return string(c.RequestCtx.Host())
}

// SetHost implements `Context#SetHost` function.
func (c *Context) SetHost(host string) {
	c.RequestCtx.Request.SetHost(host)
}

// Referer implements `Context#Referer` function.
func (c *Context) Referer() string {
	return string(c.Request.Header.Referer())
}

// ContentLength implements `Context#ContentLength` function.
func (c *Context) ContentLength() int64 {
	return int64(c.Request.Header.ContentLength())
}

// UserAgent implements `Context#UserAgent` function.
func (c *Context) UserAgent() string {
	return string(c.RequestCtx.UserAgent())
}

// RemoteAddress implements `Context#RemoteAddress` function.
func (c *Context) RemoteAddress() string {
	return c.RemoteAddr().String()
}

// RealIP implements `Context#RealIP` function.
func (c *Context) RealIP() string {
	ra := c.RemoteAddress()
	if ip := c.Request.Header.Peek(HeaderXForwardedFor); ip != nil {
		ra = string(ip)
	} else if ip := c.Request.Header.Peek(HeaderXRealIP); ip != nil {
		ra = string(ip)
	} else {
		ra, _, _ = net.SplitHostPort(ra)
	}
	return ra
}

// Method implements `Context#Method` function.
func (c *Context) Method() string {
	return string(c.RequestCtx.Method())
}

// SetMethod implements `Context#SetMethod` function.
func (c *Context) SetMethod(method string) {
	c.Request.Header.SetMethodBytes([]byte(method))
}

// URI implements `Context#URI` function.
func (c *Context) GetURI() string {
	return string(c.RequestURI())
}

// SetURI implements `Context#SetURI` function.
func (c *Context) SetURI(uri string) {
	c.Request.Header.SetRequestURI(uri)
}

// Body implements `Context#Body` function.
func (c *Context) Body() io.Reader {
	return bytes.NewBuffer(c.Request.Body())
}

// SetBody implements `Context#SetBody` function.
func (c *Context) SetBody(reader io.Reader) {
	c.Request.SetBodyStream(reader, 0)
}

// FormValue implements `Context#FormValue` function.
func (c *Context) FormValue(name string) string {
	return string(c.RequestCtx.FormValue(name))
}

// FormParams implements `Context#FormParams` function.
func (c *Context) FormParams() (params map[string][]string) {
	params = make(map[string][]string)
	mf, err := c.RequestCtx.MultipartForm()

	if err == fasthttp.ErrNoMultipartForm {
		c.PostArgs().VisitAll(func(k, v []byte) {
			key := string(k)
			if _, ok := params[key]; ok {
				params[key] = append(params[key], string(v))
			} else {
				params[string(k)] = []string{string(v)}
			}
		})
	} else if err == nil {
		for k, v := range mf.Value {
			if len(v) > 0 {
				params[k] = v
			}
		}
	}

	return
}

// FormFile implements `Context#FormFile` function.
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	return c.RequestCtx.FormFile(name)
}

// MultipartForm implements `Context#MultipartForm` function.
func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.RequestCtx.MultipartForm()
}

// Cookie implements `Context#Cookie` function.
func (ctx *Context) Cookie(name string) (*Cookie, error) {
	c := fasthttp.Cookie{}
	b := ctx.Request.Header.Cookie(name)
	if b == nil {
		return nil, ErrCookieNotFound
	}
	c.SetKey(name)
	c.SetValueBytes(b)
	return &Cookie{c}, nil
}

// Cookies implements `Context#Cookies` function.
func (ctx *Context) Cookies() []*Cookie {
	cookies := []*Cookie{}
	ctx.Request.Header.VisitAllCookie(func(name, value []byte) {
		c := fasthttp.Cookie{}
		c.SetKeyBytes(name)
		c.SetValueBytes(value)
		cookies = append(cookies, &Cookie{c})
	})
	return cookies
}
