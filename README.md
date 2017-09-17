# flick

[![](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/olafal0/flick)

Extremely simple server framework inspired by [Flask](http://flask.pocoo.org/) and [web.go](https://github.com/hoisie/web) and written in Go.

# Features
* HTTP and HTTPS serving with just `f.Serve(":80")`
* Easily define and bind handler functions with `f.Get("/route",handler)`
* Cached template rendering, even with changing data
* Easily auto-generate and use self-signed TLS: `f.ServeTLSSelfSign(":443", ":80")`
* Use websockets as just another handler function
* Automatic request logging and timing
* Additional template functions built-in, including CSS inlining

# Usage
```go
package main
import f "github.com/olafal0/flick"
func main() {
    f.Get("/", func(c *f.Context) {
        c.WriteString("Hello, world!")
    })
    f.Serve(":5000")
}
```
See the example directory for a more complete example.
