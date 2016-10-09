package flick

import (
	"bytes"
	"net/http"
	"time"
)

// Serve serves a static resource
func (c Context) Serve(filename string) {
	reader := bytes.NewReader(resource[filename].Data)
	http.ServeContent(c.wr, c.req, filename, time.Now(), reader)
}

// ServeData serves a byte array
func (c Context) ServeData(filename string, data []byte) {
	reader := bytes.NewReader(data)
	http.ServeContent(c.wr, c.req, filename, time.Now(), reader)
}

// NewPageData returns a map[string]string which can be used to pass data to templates
func NewPageData() map[string]interface{} {
	return make(map[string]interface{})
}
