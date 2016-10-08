package flick

import (
	"fmt"
	"log"
	"net/http"
)

// Serve starts the webserver
func Serve(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

// SetHandler takes a pattern string and a function(*http.Request)
// and adds it to the DefaultServeMux
func SetHandler(pattern string, handler func(w WebWriter)) {

	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Page requested: %s\n", pattern)
			handler(WebWriter{w, r})
		})

}
