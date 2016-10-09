package flick

import (
	"bytes"
	"fmt"
	"html/template"
)

var t template.Template

func init() {
	fmt.Println("Hello??")
	t := template.New("")
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
}

// RenderTemplate parses content using data and returns the resulting byte array
func RenderTemplate(filename string, data interface{}) []byte {
	fmt.Println("Creating new template")
	temp := new(template.Template)
	temp.Funcs(template.FuncMap{
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
	temp, err := temp.Parse(string(resource[filename].Data))
	if err != nil {
		fmt.Println(err)
	}
	var buf bytes.Buffer
	fmt.Println("Executing template")
	err = temp.Execute(&buf, data)
	if err != nil {
		fmt.Println(err)
	}
	return buf.Bytes()
}
