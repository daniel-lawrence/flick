package flick

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
)

// keep all the parsetrees we need pre-loaded
var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
}

func loadTemplate(filename string) *template.Template {
	// load the template file
	file, err := readFile("templates/" + filename)
	t := template.New(filename)
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
		"add": func(a, b int) int { return a + b },
		"mul": func(a, b int) int { return a * b },
	})
	t, err = t.Parse(string(file))
	if err != nil {
		fmt.Println(err)
	}
	templates[filename] = t
	return t
}

// RenderTemplate creates the template (if it doesn't exist yet) or loads it (if it does),
// renders the template, and returns the data as []byte
func RenderTemplate(filename string, data interface{}) []byte {
	var buf bytes.Buffer
	_, ok := templates[filename]
	var err error
	if !ok {
		loadTemplate(filename)
	}
	err = templates[filename].Execute(&buf, data)
	if err != nil {
		log.Println(err)
	}
	return buf.Bytes()
}

// PageData returns a map[string]string which can be used to pass data to templates
func PageData() map[string]interface{} {
	return make(map[string]interface{})
}
