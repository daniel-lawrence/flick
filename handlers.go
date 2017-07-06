package flick

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

// Get takes a pattern string and a function(*Context)
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
			c := Context{w, r, r.URL.Query()}
			handler(&c)
			elapsed := time.Since(start)
			log.Printf("%s %s: %s", r.Proto, pattern, elapsed)
		})

}

// WebSocketConnect takes a pattern string and a function(*websocket.Conn)
// that should handle a websocket connection
func WebSocketConnect(pattern string, handler func(ws *websocket.Conn)) {
	http.Handle(pattern, websocket.Handler(
		func(ws *websocket.Conn) {
			start := time.Now()
			handler(ws)
			elapsed := time.Since(start)
			log.Printf("%s: %s", pattern, elapsed)
		}))
}
