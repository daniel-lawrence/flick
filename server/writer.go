package flick

import (
	"fmt"
	"html/template"
	"net/http"
)

// WebWriter is an extended version of ReponseWriter
type WebWriter struct {
	wr http.ResponseWriter
}

// RenderTemplate takes an http.ResponseWriter, filename, and data map and renders a template
func (w WebWriter) RenderTemplate(file string, data interface{}) {
	t, err := template.ParseFiles(file)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w.wr, data)
}

// NewPageData returns a map[string]string which can be used to pass data to templates
func NewPageData() map[string]string {
	return make(map[string]string)
}
