package macross

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"mime"
	"os"
	"path"
	"path/filepath"
	"time"

	ktx "context"
)

type (

	// SerializeFunc serializes the given data of arbitrary type into a byte array.
	SerializeFunc func(data interface{}) ([]byte, error)

	// Context represents the contextual data and environment while processing an incoming HTTP request.
	Context struct {
		*fasthttp.RequestCtx
		ktx       ktx.Context   // standard context
		Serialize SerializeFunc // the function serializing the given data of arbitrary type into a byte array.
		Session   Sessioner
		macross   *Macross
		pnames    []string               // list of route parameter names
		pvalues   []string               // list of parameter values corresponding to pnames
		data      map[string]interface{} // data items managed by Get , Set , GetStore and SetStore
		index     int                    // the index of the currently executing handler in handlers
		handlers  []Handler              // the handlers associated with the current route
	}
)

const (
	indexPage = "index.html"
)

// Reset sets the request and response of the context and resets all other properties.
func (c *Context) Reset(ctx *fasthttp.RequestCtx) {
	c.RequestCtx = ctx
	c.ktx = ktx.Background()
	c.data = nil
	c.index = -1
	c.Serialize = Serialize
}

// Macross returns the Macross that is handling the incoming HTTP request.
func (c *Context) Macross() *Macross {
	return c.macross
}

func (c *Context) Kontext() ktx.Context {
	return c.ktx
}

func (c *Context) SetKontext(ktx ktx.Context) {
	c.ktx = ktx
}

func (c *Context) Handler() Handler {
	return c.handlers[c.index]
}

func (c *Context) SetHandler(h Handler) {
	c.handlers[c.index] = h
}

// Serialize converts the given data into a byte array.
// If the data is neither a byte array nor a string, it will call fmt.Sprint to convert it into a string.
func Serialize(data interface{}) (bytes []byte, err error) {
	switch data.(type) {
	case []byte:
		return data.([]byte), nil
	case string:
		return []byte(data.(string)), nil
	default:
		if data != nil {
			return []byte(fmt.Sprint(data)), nil
		}
	}
	return nil, nil
}

func (c *Context) Bind(i interface{}) error {
	return c.macross.binder.Bind(i, c)
}

func (c *Context) RequestBody() io.Reader {
	return bytes.NewBuffer(c.Request.Body())
}

// Get returns the named data item previously registered with the context by calling Set.
// If the named data item cannot be found, nil will be returned.
func (c *Context) Get(name string) interface{} {
	return c.data[name]
}

// Set stores the named data item in the context so that it can be retrieved later.
func (c *Context) Set(name string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[name] = value
}

func (c *Context) SetStore(data map[string]interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	for k, v := range data {
		c.data[k] = v
	}
}

func (c *Context) GetStore() map[string]interface{} {
	return c.data
}

// Next calls the rest of the handlers associated with the current route.
// If any of these handlers returns an error, Next will return the error and skip the following handlers.
// Next is normally used when a handler needs to do some postprocessing after the rest of the handlers
// are executed.
func (c *Context) Next() error {
	c.index++
	for n := len(c.handlers); c.index < n; c.index++ {
		if err := c.handlers[c.index](c); err != nil {
			return err
		}
	}
	return nil
}

// Abort skips the rest of the handlers associated with the current route.
// Abort is normally used when a handler handles the request normally and wants to skip the rest of the handlers.
// If a handler wants to indicate an error condition, it should simply return the error without calling Abort.
func (c *Context) Abort() error {
	c.index = len(c.handlers)
	return nil
}

// Break 中断继续执行后续动作，返回指定状态及错误，不设置错误亦可.
func (c *Context) Break(status int, err ...error) error {
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	c.Response.Header.SetStatusCode(status)
	c.macross.HandleError(c, e)
	return c.Abort()
}

// URL creates a URL using the named route and the parameter values.
// The parameters should be given in the sequence of name1, value1, name2, value2, and so on.
// If a parameter in the route is not provided a value, the parameter token will remain in the resulting URL.
// Parameter values will be properly URL encoded.
// The method returns an empty string if the URL creation fails.
func (c *Context) URL(route string, pairs ...interface{}) string {
	if r := c.macross.routes[route]; r != nil {
		return r.URL(pairs...)
	}
	return ""
}

// Data writes the given data of arbitrary type to the response.
// The method calls the Serialize() method to convert the data into a byte array and then writes
// the byte array to the response.
func (c *Context) Data(data interface{}) (err error) {
	var bytes []byte
	if bytes, err = c.Serialize(data); err == nil {
		_, err = c.Write(bytes)
	}
	return
}

func (c *Context) JSON(i interface{}, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return c.JSONBlob(b, code)
}

func (c *Context) JSONBlob(b []byte, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	return c.Blob(MIMEApplicationJSONCharsetUTF8, b, code)
}

func (c *Context) JSONP(callback string, i interface{}, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return c.JSONPBlob(callback, b, code)
}

func (c *Context) JSONPBlob(callback string, b []byte, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	c.Response.Header.Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	c.Response.Header.SetStatusCode(code)
	if _, err = c.Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.Write(b); err != nil {
		return
	}
	_, err = c.Write([]byte(");"))
	return
}

func (c *Context) Render(name string, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	if c.macross.renderer == nil {
		return ErrRendererNotRegistered
	}
	buf := new(bytes.Buffer)
	if err = c.macross.renderer.Render(buf, name, c); err != nil {
		return
	}
	c.Response.Header.Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	c.Response.Header.SetStatusCode(code)
	_, err = c.Write(buf.Bytes())
	return
}

func (c *Context) HTML(html string, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	c.Response.Header.Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	c.Response.Header.SetStatusCode(code)
	_, err = c.Write([]byte(html))
	return
}

func (c *Context) String(s string, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	c.Response.Header.Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	c.Response.Header.SetStatusCode(code)
	_, err = c.Write([]byte(s))
	return
}

func (c *Context) XML(i interface{}, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}

	b, err := xml.Marshal(i)
	if err != nil {
		return err
	}
	return c.XMLBlob(b, code)
}

func (c *Context) XMLBlob(b []byte, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	c.Response.Header.Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	c.Response.Header.SetStatusCode(code)
	if _, err = c.Write([]byte(xml.Header)); err != nil {
		return
	}
	_, err = c.Write(b)
	return
}

func (c *Context) Blob(contentType string, b []byte, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}

	c.Response.Header.Set(HeaderContentType, contentType)
	c.Response.Header.SetStatusCode(code)
	_, err = c.Write(b)
	return
}

func (c *Context) Stream(contentType string, r io.Reader, status ...int) (err error) {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	c.Response.Header.Set(HeaderContentType, contentType)
	c.Response.Header.SetStatusCode(code)
	_, err = io.Copy(c, r)
	return
}

// ServeFile serves a view file, to send a file ( zip for example) to the client
// you should use the SendFile(serverfilename,clientfilename)
//
// You can define your own "Content-Type" header also, after this function call
// This function doesn't implement resuming (by range), use ctx.SendFile/fasthttp.ServeFileUncompressed(ctx.RequestCtx,path)/fasthttpServeFile(ctx.RequestCtx,path) instead
//
// Use it when you want to serve css/js/... files to the client, for bigger files and 'force-download' use the SendFile
func (ctx *Context) ServeFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return ErrNotFound
	}
	defer f.Close()
	fi, _ := f.Stat()
	if fi.IsDir() {
		file = path.Join(file, indexPage)
		f, err = os.Open(file)
		if err != nil {
			return ErrNotFound
		}
		fi, _ = f.Stat()
	}
	return ctx.ServeContent(f, fi.Name(), fi.ModTime())
}

// SendFile sends file for force-download to the client
//
// Use this instead of ServeFile to 'force-download' bigger files to the client
func (ctx *Context) SendFile(filename string, destinationName string) {
	ctx.RequestCtx.SendFile(filename)
	ctx.RequestCtx.Response.Header.Set(HeaderContentDisposition, "attachment;filename="+destinationName)
}

func (c *Context) Attachment(file, name string) (err error) {
	return c.contentDisposition(file, name, "attachment")
}

func (c *Context) Inline(file, name string) (err error) {
	return c.contentDisposition(file, name, "inline")
}

func (c *Context) contentDisposition(file, name, dispositionType string) (err error) {
	c.Response.Header.Set(HeaderContentDisposition, fmt.Sprintf("%s; filename=%s", dispositionType, name))
	c.ServeFile(file)
	return
}

// TimeFormat is the time format to use when generating times in HTTP
// headers. It is like time.RFC1123 but hard-codes GMT as the time
// zone. The time being formatted must be in UTC for Format to
// generate the correct format.
//
// For parsing this time format, see ParseTime.
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// RequestHeader returns the request header's value
// accepts one parameter, the key of the header (string)
// returns string
func (ctx *Context) RequestHeader(k string) string {
	return string(ctx.RequestCtx.Request.Header.Peek(k))
}

// ServeContent serves content, headers are autoset
// receives three parameters, it's low-level function, instead you can use .ServeFile(string,bool)/SendFile(string,string)
//
// You can define your own "Content-Type" header also, after this function call
// Doesn't implements resuming (by range), use ctx.SendFile instead
func (ctx *Context) ServeContent(content io.ReadSeeker, filename string, modtime time.Time) error {
	if t, err := time.Parse(TimeFormat, ctx.RequestHeader(HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
		ctx.RequestCtx.Response.Header.Del(HeaderContentType)
		ctx.RequestCtx.Response.Header.Del(HeaderContentLength)
		ctx.RequestCtx.SetStatusCode(StatusNotModified)
		return nil
	}

	ctx.RequestCtx.Response.Header.Set(HeaderContentType, ctx.ContentTypeByExtension(filename))
	ctx.RequestCtx.Response.Header.Set(HeaderLastModified, modtime.UTC().Format(TimeFormat))
	ctx.RequestCtx.SetStatusCode(StatusOK)
	_, err := io.Copy(ctx.RequestCtx.Response.BodyWriter(), content)
	return err
}

// ContentTypeByExtension returns the MIME type associated with the file based on
// its extension. It returns `application/octet-stream` incase MIME type is not
// found.
func (ctx *Context) ContentTypeByExtension(name string) (t string) {
	ext := filepath.Ext(name)
	//these should be found by the windows(registry) and unix(apache) but on windows some machines have problems on this part.
	if t = mime.TypeByExtension(ext); t == "" {
		// no use of map here because we will have to lock/unlock it, by hand is better, no problem:
		if ext == ".json" {
			t = MIMEApplicationJSON
		} else if ext == ".zip" {
			t = "application/zip"
		} else if ext == ".3gp" {
			t = "video/3gpp"
		} else if ext == ".7z" {
			t = "application/x-7z-compressed"
		} else if ext == ".ace" {
			t = "application/x-ace-compressed"
		} else if ext == ".aac" {
			t = "audio/x-aac"
		} else if ext == ".ico" { // for any case
			t = "image/x-icon"
		} else if ext == ".png" {
			t = "image/png"
		} else {
			t = MIMEOctetStream
		}
	}
	return
}

func (c *Context) NoContent(status ...int) error {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusOK
	}
	c.Response.Header.SetStatusCode(code)
	return nil
}

func (c *Context) Redirect(url string, status ...int) error {
	var code int
	if len(status) > 0 {
		code = status[0]
	} else {
		code = StatusFound
	}
	if code < StatusMultipleChoices || code > StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	c.Response.Header.Set(HeaderLocation, url)
	c.Response.Header.SetStatusCode(code)
	return nil
}
