package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// TextResource is a resource to be served that doesn't change during runtime (css files for example)
type textResource struct {
	file os.FileInfo
	data []byte
}

func loadResource(filename string) *textResource {
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
	return &s
}

func (s *textResource) writeResource(output *os.File) {
	output.Write([]byte("resource[\"" + s.file.Name() + "\"] = &flick.TextResource{ Name: \"" + s.file.Name() + "\", Data: []byte(`"))
	output.Write(s.data)
	output.Write([]byte("`) }\n"))
	output.Write([]byte("resource[\"" + s.file.Name() + "\"].Modtime,_  = time.Parse(time.UnixDate,\"" + s.file.ModTime().Format(time.UnixDate) + "\")\n"))
	// output.Write([]byte("temp = resource[\"" + s.file.Name() + "\"]\n"))
	// output.Write([]byte("temp.Modtime = t\n"))
	// output.Write([]byte("resource[\"" + s.file.Name() + "\"] = temp\n"))
}
