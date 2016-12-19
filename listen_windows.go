// +build windows

// Package macross is a high productive and modular web framework in Golang.
package macross

import (
	"github.com/valyala/fasthttp"
)

func (m *Macross) Listen(args ...interface{}) {
	addr := GetAddress(args...)
	fasthttp.ListenAndServe(addr, m.HandleRequest)
}

func (m *Macross) ListenTLS(certFile, keyFile string, args ...interface{}) {
	addr := GetAddress(args...)
	fasthttp.ListenAndServeTLS(addr, certFile, keyFile, m.HandleRequest)
}

func (m *Macross) ListenTLSEmbed(certData, keyData []byte, args ...interface{}) {
	addr := GetAddress(args...)
	fasthttp.ListenAndServeTLSEmbed(addr, certData, keyData, m.HandleRequest)
}
