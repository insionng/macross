package access_test

import (
	"bytes"
	"errors"
	"fmt"
	//"github.com/stretchr/testify/assert"
	"github.com/insionng/macross"
	"github.com/insionng/macross/access"
	//"net/http"
	//"net/http/httptest"
	//"testing"
)

/*
func TestGetClientIP(t *testing.T) {
	req, _ := http.NewRequest("GET", "/users/", nil)
	req.Header.Set("X-Real-IP", "192.168.100.1")
	req.Header.Set("X-Forwarded-For", "192.168.100.2")
	req.RemoteAddr = "192.168.100.3"

	assert.Equal(t, "192.168.100.1", getClientIP(req))
	req.Header.Del("X-Real-IP")
	assert.Equal(t, "192.168.100.2", getClientIP(req))
	req.Header.Del("X-Forwarded-For")
	assert.Equal(t, "192.168.100.3", getClientIP(req))

	req.RemoteAddr = "192.168.100.3:8080"
	assert.Equal(t, "192.168.100.3", getClientIP(req))
}
*/

func getLogger(buf *bytes.Buffer) access.LogFunc {
	return func(format string, a ...interface{}) {
		fmt.Fprintf(buf, format, a...)
	}
}

func handler1(c *macross.Context) error {
	return errors.New("abc")
}
