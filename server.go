package flick

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kabukky/httpscerts"
)

// context should include some interface into all functionality
// if this is the case, then we can use a similar context type for testing
// with no difficulty stemming from http functionality

// Context is an extended version of ReponseWriter.
type Context struct {
	Wr      http.ResponseWriter
	Req     *http.Request
	Queries map[string][]string
}

// Write writes the given bytes to the reponse. This should only be called once
// per response.
func (c *Context) Write(data []byte) {
	reader := bytes.NewReader(data)
	name := c.Req.RequestURI
	http.ServeContent(c.Wr, c.Req, name, time.Now(), reader)
}

// WriteString writes the given string to the response. This should only be
// called once per response.
func (c *Context) WriteString(data string) {
	reader := bytes.NewReader([]byte(data))
	name := c.Req.RequestURI
	http.ServeContent(c.Wr, c.Req, name, time.Now(), reader)
}

func redirectToHTTPS(addr string, httpsAddr string) {
	if addr == "" {
		return
	}

	// get the host part of httpsAddr (probably the whole thing)
	httpsPort := httpsAddr[strings.Index(httpsAddr, ":"):]
	mux := http.NewServeMux()
	log.Printf("Will redirect any requests from %s to %s\n", addr, httpsPort)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := *r.URL
		url.Scheme = "https"
		url.Host = r.Host
		url.Host = url.Hostname() + ":5001"
		log.Printf("Redirecting %v to %s\n", r.Host, url.String())
		http.Redirect(w, r, url.String(), http.StatusMovedPermanently)
	})
	go http.ListenAndServe(addr, mux)
}

// Serve starts the webserver with the given address.
func Serve(addr string) {
	prepareStatics(false)
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// ServeLive starts the webserver, and reloads static files on each request.
func ServeLive(addr string) {
	prepareStatics(true)
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// ServeTLS uses a supplied certfile and keyfile to serve over HTTPS. If a redirect address is supplied,
// any requests to redirectAddress will be redirected to addr using HTTPS (for example, if redirectAddress is ":80"
// and addr is ":443", HTTP requests will be redirected to HTTPS).
func ServeTLS(addr, certfile, keyfile, redirectAddress string) {
	prepareStatics(false)
	redirectToHTTPS(redirectAddress, addr)
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, certfile, keyfile, nil))
}

// ServeTLSSelfSign auto-generates a self-signed certificate. For testing purposes only.
func ServeTLSSelfSign(addr string, redirectAddress string) {
	err := httpscerts.Check("cert.pem", "key.pem")
	//If they are not available, generate new ones.
	if err != nil {
		err := httpscerts.Generate("cert.pem", "key.pem", addr)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	prepareStatics(true)
	redirectToHTTPS(redirectAddress, addr)
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil))
}
