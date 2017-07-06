package flick

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func readFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File not found: %q\n", filename)
		return nil, err
	}
	return data, nil
}

var staticsPrepared = false

func prepareStatics(hotload bool) {
	if staticsPrepared {
		return
	}
	staticsPrepared = true
	files, err := ioutil.ReadDir("./static/")
	if err != nil {
		log.Print("No static files found")
		return
	}
	for _, f := range files {
		fmt.Println(f.Name())
		if hotload {
			serveLiveUpdating(f.Name())
		} else {
			serveStaticFile(f.Name(), f.ModTime())
		}
	}
}

func readStaticFile(filename string) (io.ReadSeeker, error) {
	path := "static/" + filename
	// get contents of file
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contentsReader := bytes.NewReader(contents)
	return contentsReader, nil
}

func serveStaticFile(filename string, modtime time.Time) {
	contentsReader, err := readStaticFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	pattern := "/" + filename
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			// make sure this is actually supposed to be a GET request
			http.ServeContent(w, r, filename, modtime, contentsReader)
			elapsed := time.Since(start)
			log.Printf("%s %s: %s", r.Proto, pattern, elapsed)
		})
}

func serveLiveUpdating(filename string) {
	pattern := "/" + filename
	http.HandleFunc(pattern,
		func(w http.ResponseWriter, r *http.Request) {
			contentsReader, err := readStaticFile(filename)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			start := time.Now()
			// make sure this is actually supposed to be a GET request
			http.ServeContent(w, r, filename, time.Now(), contentsReader)
			elapsed := time.Since(start)
			log.Printf("%s %s: %s", r.Proto, pattern, elapsed)
		})
}
