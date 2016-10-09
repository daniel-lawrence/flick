package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// textResource is a resource to be served that doesn't change during runtime (css files for example)
type textResource struct {
	file     os.FileInfo
	data     []byte
	isStatic bool
}

func loadResource(filename string, isStatic bool) *textResource {
	var s textResource
	fmt.Println(filename)
	var err error
	s.data, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Open(filename)
	defer f.Close()
	s.file, err = f.Stat()
	if err != nil {
		fmt.Println(err)
	}
	s.isStatic = isStatic
	return &s
}

func (s *textResource) writeResource(output *os.File) {
	output.Write([]byte("resource[\"" + s.file.Name() + "\"] = &flick.TextResource{ Name: \"" + s.file.Name() + "\", Data: []byte(`"))
	output.Write(s.data)
	output.Write([]byte("`), IsStatic: "))
	if s.isStatic {
		output.Write([]byte("true }\n"))
	} else {
		output.Write([]byte("false }\n"))
	}
	output.Write([]byte("resource[\"" + s.file.Name() + "\"].Modtime,_  = time.Parse(time.UnixDate,\"" + s.file.ModTime().Format(time.UnixDate) + "\")\n"))
	// output.Write([]byte("temp = resource[\"" + s.file.Name() + "\"]\n"))
	// output.Write([]byte("temp.Modtime = t\n"))
	// output.Write([]byte("resource[\"" + s.file.Name() + "\"] = temp\n"))
}
