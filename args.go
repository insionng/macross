package macross

import (
	"github.com/insionng/macross/libraries/com"
	"strconv"
	"time"
)

type (
	Args struct {
		s string
	}
)

func (a *Args) MustInt() int {
	return com.StrTo(a.s).MustInt()
}

func (a *Args) MustInt64() int64 {
	return com.StrTo(a.s).MustInt64()
}

func (a *Args) MustUint8() uint8 {
	return com.StrTo(a.s).MustUint8()
}

func (a *Args) MustUint() uint {
	return uint(com.StrTo(a.s).MustInt64())
}

func (a *Args) Float32() (f float32, e error) {
	var _f float64
	_f, e = strconv.ParseFloat(a.s, 32)
	f = float32(_f)
	return
}

func (a *Args) MustFloat32() (f float32) {
	_f, _ := strconv.ParseFloat(a.s, 32)
	f = float32(_f)
	return
}

func (a *Args) Float64() (f float64, e error) {
	f, e = strconv.ParseFloat(a.s, 64)
	return
}

func (a *Args) MustFloat64() (f float64) {
	f, _ = strconv.ParseFloat(a.s, 64)
	return
}

func (a *Args) Int() (int, error) {
	return com.StrTo(a.s).Int()
}

func (a *Args) Int64() (int64, error) {
	return com.StrTo(a.s).Int64()
}

func (a *Args) String() string {
	return com.StrTo(a.s).String()
}

func (a *Args) Bytes() []byte {
	return []byte(com.StrTo(a.s).String())
}

func (a *Args) Time() time.Time {
	tme, _ := time.Parse("2006-01-02 03:04:05 PM", com.StrTo(a.s).String())
	return tme
}

func (a *Args) Exist() bool {
	return com.StrTo(a.s).Exist()
}

func (a *Args) ToStr(args ...int) (s string) {
	return com.ToStr(a.s, args...)
}

func (a *Args) ToSnakeCase(str ...string) string {
	var s string
	if len(str) > 0 {
		s = str[0]
	} else {
		if len(a.s) != 0 {
			s = a.s
		}
	}
	return com.ToSnakeCase(s)
}

// Param returns the named parameter value that is found in the URL path matching the current route.
// If the named parameter cannot be found, an empty string will be returned.
func (c *Context) Param(name string) *Args {
	var a = new(Args)
	for i, n := range c.pnames {
		if n == name {
			a.s = c.pvalues[i]
		}
	}
	return a
}

func (c *Context) Form(key ...string) *Args {
	var a = new(Args)
	var k string
	if len(key) > 0 {
		k = key[0]
		a.s = c.FormValue(k)
	}
	return a
}

//Args 先从URL获取参数，如若没有则再尝试从from获取参数
func (c *Context) Args(key ...string) *Args {
	var a = new(Args)
	var k string
	if len(key) > 0 {
		k = key[0]
		for i, n := range c.pnames {
			if n == k {
				a.s = c.pvalues[i]
			}
		}
		if len(a.s) == 0 {
			a.s = c.FormValue(k)
		}
	}
	return a
}

func (c *Context) Parameter(i int) (value string) {
	l := len(c.pnames)
	if i < l {
		value = c.pvalues[i]
	}
	return
}

// QueryParam implements `Context#QueryParam` function.
func (c *Context) QueryParam(name string) string {
	return string(c.QueryArgs().Peek(name))
}

// QueryParams implements `Context#QueryParams` function.
func (c *Context) QueryParams() (params map[string][]string) {
	params = make(map[string][]string)
	c.URI().QueryArgs().VisitAll(func(k, v []byte) {
		_, ok := params[string(k)]
		if !ok {
			params[string(k)] = make([]string, 0)
		}
		params[string(k)] = append(params[string(k)], string(v))
	})
	return
}

func (c *Context) QueryString() string {
	return string(c.URI().QueryString())
}
