package main

import (
	f "flick/server"
	"net/http"
)

func main() {
	f.SetHandler("/", rootHandler)
	f.Serve(":5000")
}

func rootHandler(req *http.Request) string {
	d := f.NewPageData()
	d["Name"] = "Hello world!"
	d["Body"] = "This is a test."
	text := f.RenderTemplate("index.html", &d)
	return text
}
