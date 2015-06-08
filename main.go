// HttpDebugProxy is a HTTP proxy which dumps out the requests and responses.
package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// CustomTransport replaces DefaultTransport for overriding RoundTrip.
type CustomTransport struct{}

func main() {
	// TODO: improve flags.
	var src, dst string
	flag.Parse()
	args := flag.Args()
	// First arg is destination.
	if len(args) >= 1 {
		dst = args[0]
	} else {
		dst = "http://127.0.0.1:8080"
	}
	// Second arg is source.
	if len(args) == 2 {
		src = args[1]
	} else {
		src = ":80"
	}

	// Create destination url.
	u, err := url.Parse(dst)
	if err != nil {
		log.Fatal("Error parsing destination url:", err)
	}

	// Create reverse proxy.
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &CustomTransport{}

	// Create server.
	s := &http.Server{
		Addr:           src,
		Handler:        proxy,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

// RoundTrip replaces DefaultTransport RoundTrip, contains the Request and Response dumps.
func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read the Request, including body.
	reqBody, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println("Error calling DumpRequest:", err)
	}
	log.Println(">>>>> Request >>>>>")
	log.Println("\n\n" + string(reqBody))

	// Send request and receive response.
	response, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		log.Println("Error calling RoundTrip:", err)
		return nil, err
	}

	// Read the Response, including body.
	respBody, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Println("Error calling DumpResponse:", err)
		return nil, err
	}
	log.Println("<<<<< Response <<<<<")
	log.Println("\n\n" + string(respBody))

	return response, err
}
