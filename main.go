package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var (
	backend = flag.String("backend", "", "The backend url e.g. http://www.example.com")
	port    = flag.Int("p", 8080, "The port to bind")
)

func createReverseProxy(u *url.URL) http.Handler {
	r := httputil.NewSingleHostReverseProxy(u)

	// Rewrite the host name
	r.Director = func(r *http.Request) {
		r.URL.Scheme = u.Scheme
		r.URL.Host = u.Host
		r.Host = u.Host
	}

	// Rewrite redirects
	r.ModifyResponse = func(r *http.Response) error {
		loc := r.Header.Get("Location")
		if loc != "" {

			loc = strings.Replace(loc, u.String(), fmt.Sprintf("http://localhost:%v", *port), -1)
			log.Printf("%v", loc)
			r.Header.Set("Location", loc)
		}
		return nil
	}

	return r
}

func main() {
	flag.Parse()

	// Check that we have a backend
	if *backend == "" {
		log.Fatalf("No backend url specified via --backend")
	}

	// Get the URL
	u, err := url.Parse(*backend)
	if err != nil {
		log.Fatalf("Could not parse backend url '%v': %v", backend, err)
	}

	for _, arg := range flag.Args() {
		parts := strings.Split(arg, ":")
		mount := parts[0]
		path := parts[1]
		log.Printf("Mounting %v on path %v", path, mount)
		http.Handle(mount, handlers.LoggingHandler(os.Stdout, http.StripPrefix(mount, http.FileServer(http.Dir(path)))))
	}

	http.Handle("/", createReverseProxy(u))

	listen := fmt.Sprintf(":%v", *port)
	log.Printf("Forwarding %v to %v", listen, u)
	log.Fatal(http.ListenAndServe(listen, nil))
}
