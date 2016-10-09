package flick

import (
	"fmt"
	"log"
	"net/http"
)

// Context is an extended version of ReponseWriter
type Context struct {
	wr  http.ResponseWriter
	req *http.Request
}

// Serve starts the webserver
func Serve(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

// SetDefaultHandler is a generic handler that will simply serve a resource
func SetDefaultHandler(pattern string) {
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("requested: %s\n", pattern[1:])
			Context{w, r}.ServeData(pattern[1:], resource[pattern[1:]].Data)
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
