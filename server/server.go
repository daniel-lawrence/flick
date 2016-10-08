package flick

import (
	"log"
	"net/http"
)

// Serve starts the webserver
func Serve(addr string) {
	log.Fatal(http.ListenAndServe(addr, nil))
}

// SetHandler takes a pattern string and a function(http.ResponseWriter,*http.Request)
func SetHandler(pattern string, handler func(req *http.Request) string) {

	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			Write(w, handler(r))
		})

}
