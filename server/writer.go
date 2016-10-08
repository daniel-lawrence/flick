package flick

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

// WebWriter is an extended version of ReponseWriter
type WebWriter struct {
	wr http.ResponseWriter
}

// RenderTemplate takes an http.ResponseWriter, filename, and data map and renders a template
func RenderTemplate(file string, data interface{}) string {
	t, err := template.ParseFiles(file)
	if err != nil {
		fmt.Println(err)
	}
	var buf bytes.Buffer
	t.Execute(&buf, data)
	return buf.String()
}

// NewPageData returns a map[string]string which can be used to pass data to templates
func NewPageData() map[string]string {
	return make(map[string]string)
}

func Write(wr http.ResponseWriter, text string) {
	fmt.Fprintf(wr, text)
}
