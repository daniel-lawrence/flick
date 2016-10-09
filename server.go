package flick

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var resource map[string]*TextResource

// TextResource is a resource to be served that doesn't change during runtime (css files for example)
type TextResource struct {
	Modtime time.Time
	Name    string
	Data    []byte
}

// Context is an extended version of ReponseWriter
type Context struct {
	wr  http.ResponseWriter
	req *http.Request
}

// SetResources gives the flick package all the resources that the server needs
func SetResources(r map[string]*TextResource) {
	resource = r
}

// Serve starts the webserver
func Serve(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

// SetDefaultHandler is a generic handler that will simply serve a resource
func SetDefaultHandler(pattern string) {
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("requested: %s\n", pattern)
			Context{w, r}.ServeData(pattern, resource[pattern].Data)
		})
}

// SetHandler takes a pattern string and a function(*http.Request)
// and adds it to the DefaultServeMux
func SetHandler(pattern string, handler func(c Context)) {

	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Page requested: %s\n", pattern)
			handler(Context{w, r})
		})

}
