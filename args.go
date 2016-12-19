package macross

import (
	"github.com/insionng/macross/libraries/com"
	"strconv"
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

func (a *Args) Float32() (f float64, e error) {
	f, e = strconv.ParseFloat(a.s, 32)
	return
}

func (a *Args) Float64() (f float64, e error) {
	f, e = strconv.ParseFloat(a.s, 64)
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
