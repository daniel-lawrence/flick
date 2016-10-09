package flick

import (
	"bytes"
	"net/http"
)

// Serve serves a static resource
func (c Context) Serve(filename string) {
	reader := bytes.NewReader(resource[filename].Data)
	http.ServeContent(c.wr, c.req, filename, resource[filename].Modtime, reader)
}

// ServeData serves a byte array
func (c Context) ServeData(filename string, data []byte) {
	reader := bytes.NewReader(data)
	http.ServeContent(c.wr, c.req, filename, resource[filename].Modtime, reader)
}
