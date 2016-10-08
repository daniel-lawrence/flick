package main

import (
	f "flick"
	"net/http"
)

func main() {
	f.SetHandler("/", rootHandler)
	f.Serve(":5000")
}

func rootHandler(wr f.WebWriter, req *http.Request) {
	d := f.NewPageData()
	d["Name"] = "Hello world!"
	d["Body"] = "This is a test."
	wr.RenderTemplate("index.html", &d)
}
