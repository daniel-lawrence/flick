package flick

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
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
		log.Print(err)
	}
	fmt.Print("Adding static files:\n")
	for _, f := range files {
		fmt.Println(f.Name())
		serveStaticFile(f.Name(), f.ModTime())
	}
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Get takes a pattern string and a function(*http.Request)
// and adds it to the DefaultServeMux
func Get(pattern string, handler func(c *Context)) {

	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// make sure this is actually supposed to be a GET request
			if r.Method != "" && r.Method != "GET" {
				// use reflection to get the name of the handler method
				methodName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
				log.Printf("Warning: GET handler for function %s got non-GET method type", methodName)
			}
			handler(&Context{w, r})
			elapsed := time.Since(start)
			log.Printf("%s %s: %s", r.Proto, pattern, elapsed)
		})

}

func serveStaticFile(filename string, modtime time.Time) {
	path := "static/" + filename
	// get contents of file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error serving %s: %v", path, err)
		return
	}
	contentsReader := bytes.NewReader(contents)
	Get("/"+filename,
		func(c *Context) {
			http.ServeContent(c.Wr, c.Req, filename, modtime, contentsReader)
		})
}
