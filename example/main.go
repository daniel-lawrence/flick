package main

import (
	"net/http"
	"time"

	f "github.com/olafal0/flick"
)

func main() {
	// assign a handler function for the root path
	f.Get("/", rootHandler)
	// example of a handler returning a file
	f.Get("/test", staticHandler)
	// start serving on 0.0.0.0:5000
	// do it in a goroutine so we can serve HTTP and HTTPS simulataneously
	go f.Serve(":5000")
	// we can also generate our own certs and serve over HTTPS
	// this will generate security warnings in browsers, though
	// (it's also bad in general - only use it for testing.
	//  ServeTLS() is what you really want)
	f.ServeTLSSelfSign(":5001", "")
	// visit http://localhost:5000 or https://localhost:5001

	// if you wanted to automatically redirect HTTP on :5000 to HTTPS on :5001,
	// you could just do this:
	// f.ServeTLSSelfSign(":5001", ":5000")
}

// A handler just takes a context. Use f.Context.Write([]byte) for your response.
func rootHandler(c *f.Context) {
	// 404 if it isn't the correct path. Because of the way Go's serve mux works
	// by default, "/" acts as a catch-all for any path that doesn't exist.
	if c.Req.URL.Path != "/" {
		http.NotFound(c.Wr, c.Req)
		return
	}
	// make the data to pass to the template
	// this could easily be replaced with simply passing time.Now() to the renderer,
	// but this is an example of something that could be extended.
	// In the template, you can just use {{.ServerTime}} to print the value
	data := make(map[string]interface{})
	data["ServerTime"] = time.Now()
	// data could also be a struct with a field called ServerTime.

	// Write the rendered template.
	// The first time this handler is called, this will take maybe 10ms.
	// Most of that is overhead from reading the file.
	// After that, the template is cached - calling the function again, even
	// with different data, only takes around 100µs (1/100th of the original time).
	c.Write(f.RenderTemplate("index.html", data))
}

func staticHandler(c *f.Context) {
	// Normally, there are two types of serving:
	// - just serving static files
	// - serving templates with contents that change based on context
	// But if a route (e.g. /admin) needs to return something different
	// based on context, templates would be extremely clunky.
	// So, there's a way to programmatically return a static file.
	c.Write(f.RenderStatic("test.txt", true))
}
