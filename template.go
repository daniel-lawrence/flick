package flick

import (
	"bytes"
	"fmt"
	"html/template"
)

// keep all the parsetrees we need pre-loaded
var templates map[string]*template.Template

func preloadTemplates() {
	templates = make(map[string]*template.Template)
	for k, v := range resource {
		templates[k] = template.New(k)
		templates[k].Funcs(template.FuncMap{
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
		var err error
		templates[k], err = templates[k].Parse(string(v.Data))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// RenderTemplate parses content using data and returns the resulting byte array
func RenderTemplate(filename string, data interface{}) []byte {
	var buf bytes.Buffer
	fmt.Println("Executing template")
	err := templates[filename].Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}

// NewPageData returns a map[string]string which can be used to pass data to templates
func NewPageData() map[string]interface{} {
	return make(map[string]interface{})
}
