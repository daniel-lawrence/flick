package flick

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Context is an extended version of ReponseWriter
type Context struct {
	Wr  http.ResponseWriter
	Req *http.Request
}

func (c *Context) Write(data []byte) {
	reader := bytes.NewReader(data)
	name := c.Req.RequestURI
	http.ServeContent(c.Wr, c.Req, name, time.Now(), reader)
}

// Serve starts the webserver
func Serve(addr string) {
	// add handlers for static files automatically
	files, err := ioutil.ReadDir("./static/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Printf("Adding static file to list: %s\n", f.Name())
		serveStaticFile(f.Name())
	}
	fmt.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Get takes a pattern string and a function(*http.Request)
// and adds it to the DefaultServeMux
func Get(pattern string, handler func(c *Context)) {

	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			handler(&Context{w, r})
			elapsed := time.Since(start)
			log.Printf("%s %s: %s", r.Proto, pattern, elapsed)
		})

}

func serveStaticFile(filename string) {
	path := "static/" + filename
	Get("/"+filename,
		func(c *Context) {
			http.ServeFile(c.Wr, c.Req, path)
		})
}
