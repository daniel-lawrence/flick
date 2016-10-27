package flick

import (
	"fmt"
	"io/ioutil"
)

func readFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File not found: %q\n", filename)
		return nil, err
	}
	return data, nil
}
