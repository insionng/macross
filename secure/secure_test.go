package secure_test

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/secure"
	"testing"
)

func TestSecure(t *testing.T) {
	e := macross.New()
	e.Use(secure.Secure())
	go e.Run(":8000")

	m := macross.New()
	m.Use(secure.SecureWithConfig(secure.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	go m.Run(":9000")

}
