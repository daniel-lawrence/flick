package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// PackFiles packs all files into go files to be included in the binary
func main() {
	basePath := os.Args[1]

	out, err := os.Create(basePath + "/files.go")
	if err != nil {
		fmt.Println(err)
	}
	out.Write([]byte("package main \n\nimport \"github.com/olafal0/flick\"\nvar resource map[string]*flick.TextResource\n\nfunc setMap () {\n"))
	out.Write([]byte("resource = make(map[string]*flick.TextResource)\n"))

	sfs, err := ioutil.ReadDir(basePath + "/static/")
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range sfs {
		s := loadResource(basePath+"/static/"+f.Name(), true)
		s.writeResource(out)
	}

	tfs, err := ioutil.ReadDir(basePath + "/templates/")
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range tfs {
		s := loadResource(basePath+"/templates/"+f.Name(), false)
		s.writeResource(out)
	}
	out.Write([]byte("}\n"))
	out.Close()
}
