package flick

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"time"
)

// StaticResource represents a resource that doesn't change during runtime
type StaticResource struct {
	data    []byte
	modtime time.Time
	Name    string
}

// LoadStatic loads a file
func LoadStatic(file string) (s StaticResource) {
	s.Name = file
	fileData, fError := os.Stat(file)
	if fError != nil {
		// 404
		fmt.Println("File error")
	} else {
		// set ETag header
		s.modtime = fileData.ModTime()
	}
	s.data, fError = ioutil.ReadFile(file)
	if fError != nil {
		// 404
		fmt.Println("File error")
	}
	return s
}

// LoadFromTemplate loads a template with a static set of data
func LoadFromTemplate(file string, data interface{}) (s StaticResource) {
	s.Name = file
	fileData, fError := os.Stat(file)
	if fError != nil {
		// 404
		fmt.Println("File error")
	} else {
		// set ETag header
		s.modtime = fileData.ModTime()
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
	s.data = buf.Bytes()
	return s
}
