// +build linux darwin dragonfly freebsd netbsd openbsd rumprun

// Package macross is a high productive and modular web framework in Golang.
package macross

import (
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"runtime"
)

func (m *Macross) Listen(args ...interface{}) {
	addr := GetAddress(args...)
	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
		ln, err := reuseport.Listen("tcp4", addr)
		if err != nil {
			log.Fatalf("error in reuseport.Listen: %s", err)
		}

		if err = fasthttp.Serve(ln, m.ServeHTTP); err != nil {
			log.Fatalf("error in fasthttp.Serve: %s", err)
		}
	} else {
		fasthttp.ListenAndServe(addr, m.ServeHTTP)
	}
}

func (m *Macross) ListenTLS(certFile, keyFile string, args ...interface{}) {
	addr := GetAddress(args...)
	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
		ln, err := reuseport.Listen("tcp4", addr)
		if err != nil {
			log.Fatalf("error in reuseport.Listen: %s", err)
		}

		if err = fasthttp.ServeTLS(ln, certFile, keyFile, m.ServeHTTP); err != nil {
			log.Fatalf("error in fasthttp.ServeTLS: %s", err)
		}

	} else {
		fasthttp.ListenAndServeTLS(addr, certFile, keyFile, m.ServeHTTP)
	}
}

func (m *Macross) ListenTLSEmbed(certData, keyData []byte, args ...interface{}) {
	addr := GetAddress(args...)
	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(runtime.NumCPU())
		ln, err := reuseport.Listen("tcp4", addr)
		if err != nil {
			log.Fatalf("error in reuseport.Listen: %s", err)
		}

		if err = fasthttp.ServeTLSEmbed(ln, certData, keyData, m.ServeHTTP); err != nil {
			log.Fatalf("error in fasthttp.ServeTLSEmbed: %s", err)
		}

	} else {
		fasthttp.ListenAndServeTLSEmbed(addr, certData, keyData, m.ServeHTTP)
	}
}
