package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/rbxb/flashboard"
	"github.com/rbxb/httpfilter"
)

var port string
var root string
var cap int
var postSize int
var lifetime int

func init() {
	flag.StringVar(&port, "port", ":8080", "The address and port the fileserver listens at.")
	flag.StringVar(&root, "root", "./root", "The directory serving files.")
	flag.IntVar(&cap, "cap", 20, "The maximum number of posts.")
	flag.IntVar(&postSize, "postSize", 200, "The maximum post length.")
	flag.IntVar(&lifetime, "lifetime", 30, "The lifetime of a post in seconds.")
}

func main() {
	flag.Parse()
	flashboardSv := flashboard.NewServer(postSize, cap, time.Second*time.Duration(lifetime))
	fs := httpfilter.NewServer(root, "", map[string]httpfilter.OpFunc{
		"flashboard": func(w http.ResponseWriter, req *http.Request, args ...string) {
			flashboardSv.ServeHTTP(w, req)
		},
	})
	server := http.Server{
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)), //disable HTTP/2
		Addr:         port,
		Handler:      fs,
	}
	log.Fatal(server.ListenAndServe())
}
