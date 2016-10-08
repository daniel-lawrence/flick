package flick

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

// WebWriter is an extended version of ReponseWriter
type WebWriter struct {
	wr  http.ResponseWriter
	req *http.Request
}

// RenderTemplate takes a filename and data map and renders a template
func (w WebWriter) RenderTemplate(file string, data interface{}) {
	// get file metadata
	fileData, fError := os.Stat(file)
	if fError != nil {
		// 404
		fmt.Println("File error")
	} else {
		// set ETag header
		w.wr.Header().Set("Etag", fileData.ModTime().String())
	}
	if match := w.req.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, fileData.ModTime().String()) {
			w.wr.WriteHeader(http.StatusNotModified)
			return
		}
	}
	t := template.New(file)
	t.Funcs(template.FuncMap{
		"mod": func(i, j int) int { return i % j },
		"loop": func(n int) []int {
			loop := make([]int, n)
			for i := 0; i < n; i++ {
				loop[i] = i
			}
			return loop
		},
		"sub": func(a, b int) int { return a - b },
		"div": func(a, b int) int { return a / b },
	})
	t, err := t.ParseFiles(file)
	if t != nil {
		fmt.Println(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
	}
	reader := bytes.NewReader(buf.Bytes())
	http.ServeContent(w.wr, w.req, file, fileData.ModTime(), reader)
}

// Serve serves a filename without template parsing
func (w WebWriter) Serve(file string) {
	// get file metadata
	fileData, fError := os.Stat(file)
	if fError != nil {
		// 404
		fmt.Println("File error")
	} else {
		// set ETag header
		w.wr.Header().Set("Etag", fileData.ModTime().String())
	}
	http.ServeFile(w.wr, w.req, file)
}

// ServeStatic serves a static resource
func (w WebWriter) ServeStatic(s StaticResource) {
	reader := bytes.NewReader(s.data)
	http.ServeContent(w.wr, w.req, s.Name, s.modtime, reader)
}

// NewPageData returns a map[string]string which can be used to pass data to templates
func NewPageData() map[string]interface{} {
	return make(map[string]interface{})
}
