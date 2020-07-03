package main

import (
	"flag"
	"inbox"
	"log"
	"net/http"
	"time"
)

func main() {
	bind := flag.String("bind", "0.0.0.0:3000", "a bind for the http server")
	serverCert := flag.String("server-cert", "server.crt", "HTTPS server certificate file")
	serverKey := flag.String("server-key", "server.key", "HTTPS server private key file")
	cleanupInterval := flag.Int("cleanup-interval", 1, "Interval in seconds between server cleaning up inboxes")
	public := flag.String("public", "public", "Directory path of static files to serve")
	https := flag.Bool("https", true, "Run server in HTTPS mode or HTTP")
	cors := flag.Bool("cors", false, "Allow CORS")

	flag.Parse()

	server := inbox.Server{
		CORS:            *cors,
		Mailboxes:       inbox.New(),
		CleanupInterval: time.Second * time.Duration(*cleanupInterval),
	}
	go server.Clean()

	http.Handle("/", http.FileServer(http.Dir(*public)))
	http.Handle("/inbox", &server)

	if *https {
		log.Fatal(http.ListenAndServeTLS(*bind, *serverCert, *serverKey, nil))
	} else {
		log.Fatal(http.ListenAndServe(*bind, nil))
	}
}
