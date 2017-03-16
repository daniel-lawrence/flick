package flick

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
)

// keep all the parsetrees we need pre-loaded
var templates map[string]*template.Template

// keep static resources cached too when they are used explicitly
var staticCached map[string]*[]byte

func init() {
	templates = make(map[string]*template.Template)
	staticCached = make(map[string]*[]byte)
}

func loadTemplate(filename string) *template.Template {
	// load the template file
	file, _ := readFile("templates/" + filename)
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
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"div": func(a, b int) int { return a / b },
		"mul": func(a, b int) int { return a * b },
		"inlinecss": func(filename string) template.CSS {
			return template.CSS(RenderStatic(filename, true))
		},
	})
	t, err := t.Parse(string(file))
	if err != nil {
		fmt.Println("Warning: template parsing error! May not be injection safe!")
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

// RenderStatic renders a static file, and returns its data
func RenderStatic(filename string, cache bool) []byte {
	if cache {
		_, ok := staticCached[filename]
		if ok {
			return *staticCached[filename]
		}
	}
	path := "static/" + filename
	// get contents of file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error serving %s: %v", path, err)
		return nil
	}
	if cache {
		staticCached[filename] = &contents
	}
	return contents
}
