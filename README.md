# Macross

The Macross Web Framework By Insion

## Requirements

Go 1.7.4 or above.

## Installation

Run the following command to install the package:

```
go get -u github.com/insionng/macross
```

## Getting Started

Create a `server.go` file with the following content:

```go
package main

import (
	"github.com/insionng/macross"
)

func main() {
	m := macross.New()
	
	m.Get("/", func(self *macross.Context) error {
		return self.String("Hello, Macross")
	})

	m.Listen(9000)
}
```

Now run the following command to start the Web server:

```
go run server.go
```

You should be able to access URLs such as `http://localhost:9000`.


## Getting Started via JWT

```go
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
```


## Getting Started via Session

```go
package main

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/recover"
	"github.com/macross-contrib/session"
	_ "github.com/macross-contrib/session/redis"
	"log"
)

func main() {

	v := macross.New()
	v.Use(recover.Recover())
	//v.Use(session.Sessioner(session.Options{"file", `{"cookieName":"MacrossSessionId","gcLifetime":3600,"providerConfig":"./data/session"}`}))
	v.Use(session.Sessioner(session.Options{"redis", `{"cookieName":"MacrossSessionId","gcLifetime":3600,"providerConfig":"127.0.0.1:6379"}`}))

	v.Get("/get", func(self *macross.Context) error {
		value := "nil"
		valueIf := self.Session.Get("key")
		if valueIf != nil {
			value = valueIf.(string)
		}

		return self.String(value)

	})

	v.Get("/set", func(self *macross.Context) error {

		val := self.QueryParam("v")
		if len(val) == 0 {
			val = "value"
		}

		err := self.Session.Set("key", val)
		if err != nil {
			log.Printf("sess.set %v \n", err)
		}
		return self.String("okay")
	})

	v.Listen(7777)
}

```

## Getting Started via i18n

```go
package main

import (
	"fmt"
	"github.com/insionng/macross"
	"github.com/macross-contrib/i18n"
)

func main() {
	m := macross.Classic()
	m.Use(i18n.I18n(i18n.Options{
		Directory:   "locale",
		DefaultLang: "zh-CN",
		Langs:       []string{"en-US", "zh-CN"},
		Names:       []string{"English", "简体中文"},
		Redirect:    true,
	}))

	m.Get("/", func(self *macross.Context) error {
		fmt.Println("Header>", self.Request.Header.String())
		return self.String("current language is " + self.Language())
	})

	// Use in handler.
	m.Get("/trans", func(self *macross.Context) error {
		fmt.Println("Header>", self.Request.Header.String())
		return self.String(fmt.Sprintf("hello %s", self.Tr("world")))
	})

	fmt.Println("Listen on 9999")
	m.Listen(9999)
}

```


## Getting Started via Go template

```go
package main

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/gonder"
	"github.com/insionng/macross/logger"
	"github.com/insionng/macross/recover"
	"github.com/insionng/macross/static"
)

func main() {
	v := macross.New()
	v.Use(logger.Logger())
	v.Use(recover.Recover())
	v.SetRenderer(gonder.Renderor())
	v.Use(static.Static("static"))
	v.Get("/", func(self *macross.Context) error {
		var data = make(map[string]interface{})
		data["name"] = "Insion Ng"
		self.SetStore(data)

		self.SetStore(map[string]interface{}{
			"title": "你好，世界",
			"oh":    "no",
		})
		self.Set("oh", "yes") //覆盖前面指定KEY
		return self.Render("index")
	})

	v.Listen(":9000")
}

```

templates/index.html
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="/static/index.js" charset="utf-8"></script>
<title>{{ .title }}</title>
</head>
<body>
    <p>{{ .oh }}</p
    <p>{{ .name }}</p>
</body>
</html>

```



## Getting Started via Pongo template

```go
package main

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/logger"
	"github.com/insionng/macross/pongor"
	"github.com/insionng/macross/recover"
	"github.com/insionng/macross/static"
)

func main() {
	v := macross.New()
	v.Use(logger.Logger())
	v.Use(recover.Recover())
	v.SetRenderer(pongor.Renderor())
	v.Use(static.Static("static"))
	v.Get("/", func(self *macross.Context) error {
		var data = make(map[string]interface{})
		data["name"] = "Insion Ng"
		self.SetStore(data)

		self.SetStore(map[string]interface{}{
			"title": "你好，世界",
			"oh":    "no",
		})
		self.Set("oh", "yes") //覆盖前面指定KEY
		return self.Render("index")
	})

	v.Listen(":9000")
}

```

templates/index.html
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="/static/index.js" charset="utf-8"></script>
<title>{{ title }}</title>
</head>
<body>
    <p>{{ oh }}</p
    <p>{{ name }}</p>
</body>
</html>

```

## Getting Started via FastTemplate

```go
package main

import (
	"github.com/insionng/macross"
	"github.com/insionng/macross/fempla"
	"github.com/insionng/macross/logger"
	"github.com/insionng/macross/recover"
	"github.com/insionng/macross/static"
)

func main() {

	v := macross.New()
	v.Use(logger.Logger())
	v.Use(recover.Recover())
	v.SetRenderer(fempla.Renderor())
	v.Use(static.Static("static"))
	v.Get("/", func(self *macross.Context) error {
		data := make(map[string]interface{})
		data["oh"] = "no"
		data["name"] = "Insion Ng"
		self.Set("title", "你好，世界")
		self.SetStore(data)
		self.Set("oh", "yes")
		return self.Render("index")
	})

	v.Listen(":9000")

}

```

templates/index.html
```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="/static/index.js" charset="utf-8"></script>
<title>{{title}}</title>
</head>
<body>
    <p>
        {{oh}}
    </p>
    <p>
        {{name}}
    </p>
</body>
</html>

```


### Routes

macross works by building a macross table in a macross and then dispatching HTTP requests to the matching handlers 
found in the macross table. An intuitive illustration of a macross table is as follows:


Routes              |  Handlers
--------------------|-----------------
`GET /users`        |  m1, m2, h1, ...
`POST /users`       |  m1, m2, h2, ...
`PUT /users/<id>`   |  m1, m2, h3, ...
`DELETE /users/<id>`|  m1, m2, h4, ...


For an incoming request `GET /users`, the first route would match and the handlers m1, m2, and h1 would be executed.
If the request is `PUT /users/123`, the third route would match and the corresponding handlers would be executed.
Note that the token `<id>` can match any number of non-slash characters and the matching part can be accessed as 
a path parameter value in the handlers.

**If an incoming request matches multiple routes in the table, the route added first to the table will take precedence.
All other matching routes will be ignored.**

The actual implementation of the macross table uses a variant of the radix tree data structure, which makes the macross
process as fast as working with a hash table.

To add a new route and its handlers to the macross table, call the `To` method like the following:
  
```go
m := macross.New()
m.To("GET", "/users", m1, m2, h1)
m.To("POST", "/users", m1, m2, h2)
```

You can also use shortcut methods, such as `Get`, `Post`, `Put`, etc., which are named after the HTTP method names:
 
```go
m.Get("/users", m1, m2, h1)
m.Post("/users", m1, m2, h2)
```

If you have multiple routes with the same URL path but different HTTP methods, like the above example, you can 
chain them together as follows,

```go
m.Get("/users", m1, m2, h1).Post(m1, m2, h2)
```

If you want to use the same set of handlers to handle the same URL path but different HTTP methods, you can take
the following shortcut:

```go
m.To("GET,POST", "/users", m1, m2, h)
```

A route may contain parameter tokens which are in the format of `<name:pattern>`, where `name` stands for the parameter
name, and `pattern` is a regular expression which the parameter value should match. A token `<name>` is equivalent
to `<name:[^/]*>`, i.e., it matches any number of non-slash characters. At the end of a route, an asterisk character
can be used to match any number of arbitrary characters. Below are some examples:

* `/users/<username>`: matches `/users/root`
* `/users/accnt-<id:\d+>`: matches `/users/accnt-123`, but not `/users/accnt-root`
* `/users/<username>/*`: matches `/users/root/profile/address`

When a URL path matches a route, the matching parameters on the URL path can be accessed via `Context.Param()`:

```go
m := macross.New()

m.Get("/users/<username>", func (self *macross.Context) error {
	username := self.Param("username").String()
	s := fmt.Sprintf("Name: %s", username)
	return self.String(s)
})
```


### Route Groups

Route group is a way of grouping together the routes which have the same route prefix. The routes in a group also
share the same handlers that are registered with the group via its `Use` method. For example,

```go
m := macross.New()
api := m.Group("/api")
api.Use(m1, m2)
api.Get("/users", h1).Post(h2)
api.Put("/users/<id>", h3).Delete(h4)
```

The above `/api` route group establishes the following macross table:


Routes                  |  Handlers
------------------------|-------------
`GET /api/users`        |  m1, m2, h1, ...
`POST /api/users`       |  m1, m2, h2, ...
`PUT /api/users/<id>`   |  m1, m2, h3, ...
`DELETE /api/users/<id>`|  m1, m2, h4, ...


As you can see, all these routes have the same route prefix `/api` and the handlers `m1` and `m2`. In other similar
macross frameworks, the handlers registered with a route group are also called *middlewares*.

Route groups can be nested. That is, a route group can create a child group by calling the `Group()` method. The macross
serves as the top level route group. A child group inherits the handlers registered with its parent group. For example, 

```go
m := macross.New()
m.Use(m1)

api := m.Group("/api")
api.Use(m2)

users := group.Group("/users")
users.Use(m3)
users.Put("/<id>", h1)
```

Because the macross serves as the parent of the `api` group which is the parent of the `users` group, 
the `PUT /api/users/<id>` route is associated with the handlers `m1`, `m2`, `m3`, and `h1`.


### Router

Router manages the macross table and dispatches incoming requests to appropriate handlers. A macross instance is created
by calling the `macross.New()` method.

To hook up macross with fasthttp, use the following code:

```go
m := macross.New()
m.Listen(":9000") 
```


### Handlers

A handler is a function with the signature `func(*macross.Context) error`. A handler is executed by the macross if
the incoming request URL path matches the route that the handler is associated with. Through the `macross.Context` 
parameter, you can access the request information in handlers.

A route may be associated with multiple handlers. These handlers will be executed in the order that they are registered
to the route. The execution sequence can be terminated in the middle using one of the following two methods:

* A handler returns an error: the macross will skip the rest of the handlers and handle the returned error.
* A handler calls `Context.Abort()` and `Context.Break()`: the macross will simply skip the rest of the handlers. There is no error to be handled.
 
A handler can call `Context.Next()` to explicitly execute the rest of the unexecuted handlers and take actions after
they finish execution. For example, a response compression handler may start the output buffer, call `Context.Next()`,
and then compress and send the output to response.


### Context

For each incoming request, a `macross.Context` object is passed through the relevant handlers. Because `macross.Context`
embeds `fasthttp.RequestCtx`, you can access all properties and methods provided by the latter.
 
Additionally, the `Context.Param()` method allows handlers to access the URL path parameters that match the current route.
Using `Context.Get()` , `Context.GetStore()` and `Context.Set()` , `Context.SetStore()`, handlers can share data between each other. For example, an authentication
handler can store the authenticated user identity by calling `Context.Set()`, and other handlers can retrieve back
the identity information by calling `Context.Get()`.

Context also provides a handy `Data()` method that can be used to write data of arbitrary type to the response.
The `Data()` method can also be overridden (by replacement) to achieve more versatile response data writing. 


### Error Handling

A handler may return an error indicating some erroneous condition. Sometimes, a handler or the code it calls may cause
a panic. Both should be handled properly to ensure best user experience. It is recommended that you use 
the `fault.Recover` handler or a similar error handler to handle these errors.

If an error is not handled by any handler, the macross will handle it by calling its `HandleError()` method which
simply sets an appropriate HTTP status code and writes the error message to the response.

When an incoming request has no matching route, the macross will call the handlers registered via the `Router.NotFound()`
method. All the handlers registered via `Router.Use()` will also be called in advance. By default, the following two
handlers are registered with `Router.NotFound()`:

* `macross.MethodNotAllowedHandler`: a handler that sends an `Allow` HTTP header indicating the allowed HTTP methods for a requested URL
* `macross.NotFoundHandler`: a handler triggering 404 HTTP error


### Middleware

* [Cache](https://github.com/macross-contrib/cache): `Middleware cache provides cache management for Macross. It can use many cache adapters, including memory, file, Redis.`
* [Session](https://github.com/macross-contrib/session): `The session package is a Macross session manager. It can use many session providers.`
* [I18n](https://github.com/macross-contrib/i18n): `Middleware i18n provides app Internationalization and Localization for Macross.`
* [Macrus](https://github.com/macross-contrib/macrus): `Package macrus provides a middleware for macross that logs request details via the logrus logging library.`
* [Statio](https://github.com/macross-contrib/statio): `serves static files and http file system.`
* [Captcha](https://github.com/macross-contrib/captcha): `Middleware captcha provides captcha service for Macross).`
* [Macrof](https://github.com/macross-contrib/macrof): `A wrapper for macross to use net/http/pprof easily.`


### Contributes

Thanks to the fasthttp, com, echo/vodka, iris, gin, beego, fasthttp-routing, FastTemplate, Pongo2, Jwt-go. And all other Go package dependencies projects


### Recipes

- [Yougam](http://www.yougam.com/) the discuss system project
- [Vuejsto](https://github.com/macross-contrib/vuejsto) the vue.js + macross example project
- [Drupapp](https://github.com/insionng/drupapp) the drupal models + macross project