package flick

import (
	"bytes"
	"log"
	"net/http"
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

func (c *Context) Write(data []byte) {
	reader := bytes.NewReader(data)
	name := c.Req.RequestURI
	http.ServeContent(c.Wr, c.Req, name, time.Now(), reader)
}

func (c *Context) WriteString(data string) {
	reader := bytes.NewReader([]byte(data))
	name := c.Req.RequestURI
	http.ServeContent(c.Wr, c.Req, name, time.Now(), reader)
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

// ServeTLS uses a supplied certfile and keyfile to serve over HTTPS.
func ServeTLS(addr, certfile, keyfile string) {
	prepareStatics(false)
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, certfile, keyfile, nil))
}

// ServeTLSSelfSign auto-generates a self-signed certificate. For testing purposes only.
func ServeTLSSelfSign(addr string) {
	err := httpscerts.Check("cert.pem", "key.pem")
	//If they are not available, generate new ones.
	if err != nil {
		err := httpscerts.Generate("cert.pem", "key.pem", addr)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	prepareStatics(false)
	log.Printf("Serving on %s\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, "cert.pem", "key.pem", nil))
}
