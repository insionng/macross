package main

import (
	"fmt"
	"github.com/insionng/macross"
	"github.com/insionng/macross/cors"
	"github.com/insionng/macross/jwt"
	"github.com/insionng/macross/logger"
	"github.com/insionng/macross/recover"
	"time"
)

/*
curl -I -X GET http://localhost:9000/jwt/get/ -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTQ3OTQ4NDUzOH0.amQOtO0GESwLoevaGSoR55jCUqZ6vsIi9DPTkDh4tSk"
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0    26    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0HTTP/1.1 200 OK
Server: Macross
Date: Fri, 18 Nov 2016 15:55:18 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 26
Vary: Origin
Access-Control-Allow-Origin: *
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTQ3OTQ4NDU3OH0.KBTm7A3xqWmQ6NLfUecfowgszfKzwMrjO3k0gf8llc8
*/

func main() {
	m := macross.New()
	m.Use(logger.Logger())
	m.Use(recover.Recover())
	m.Use(cors.CORS())

	m.Get("/", func(self *macross.Context) error {
		fmt.Println(self.Response.Header.String())
		var data = map[string]interface{}{}
		data["version"] = "1.0.0"
		return self.JSON(data)
	})

	var secret = "secret"
	var exprires = time.Minute * 1
	// 给用户返回token之前请先密码验证用户身份
	m.Post("/signin/", func(self *macross.Context) error {

		fmt.Println(self.Response.String())

		username := self.Args("username").String()
		password := self.Args("password").String()
		if (username == "insion") && (password == "PaSsworD") {
			claims := jwt.NewMapClaims()
			claims["UserId"] = 1
			claims["exp"] = time.Now().Add(exprires).Unix()

			tk, _ := jwt.NewTokenString(secret, "HS256", claims)

			var data = map[string]interface{}{}
			data["token"] = tk

			return self.JSON(data)
		}

		herr := new(macross.HTTPError)
		herr.Message = "ErrUnauthorized"
		herr.Status = macross.StatusUnauthorized
		return self.JSON(herr, macross.StatusUnauthorized)

	})

	g := m.Group("/jwt", jwt.JWT(secret))
	// http://localhost:9000/jwt/get/
	g.Get("/get/", func(self *macross.Context) error {

		var data = map[string]interface{}{}

		claims := jwt.GetMapClaims(self)
		jwtUserId := claims["UserId"].(float64)
		fmt.Println(jwtUserId)
		exp := int64(claims["exp"].(float64))
		exptime := time.Unix(exp, 0).Sub(time.Now())

		if (exptime > 0) && (exptime < (exprires / 3)) {
			fmt.Println("exptime will be expires")
			claims["UserId"] = 1
			claims["exp"] = time.Now().Add(exprires).Unix()

			token := jwt.NewToken("HS256", claims)
			tokenString, _ := token.SignedString([]byte(secret))

			self.Response.Header.Set(macross.HeaderAccessControlExposeHeaders, "Authorization")
			self.Response.Header.Set("Authorization", jwt.Bearer+" "+tokenString)
			self.Set(jwt.DefaultJWTConfig.ContextKey, token)
		}

		data["value"] = "Hello, Macross"
		return self.JSON(data)
	})

	m.Listen(":9000")
}
