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
	inboxTimeout := flag.Int("inbox-timeout", 60, "Number of seconds for the inbox to be inactive before deleting")
	messageTimeout := flag.Int("message-timeout", 60, "Number of seconds for the message to be saved in the inbox before deleting")
	public := flag.String("public", "public", "Directory path of static files to serve")
	https := flag.Bool("https", true, "Run server in HTTPS mode or HTTP")
	cors := flag.Bool("cors", false, "Allow CORS")
	maxBodySize := flag.Int64("max-body-size", 1*1024*1024, "Maximum request body size in bytes")
	maxHeaderSize := flag.Int("max-header-size", http.DefaultMaxHeaderBytes, "Maximum request body size in bytes")

	flag.Parse()

	mailboxes := inbox.New()
	mailboxes.InboxTimeout = time.Second * time.Duration(*inboxTimeout)
	mailboxes.MessageTimeout = time.Second * time.Duration(*messageTimeout)

	server := inbox.Server{
		CORS:            *cors,
		Mailboxes:       mailboxes,
		CleanupInterval: time.Second * time.Duration(*cleanupInterval),
		MaxBodySize:     *maxBodySize,
	}

	go server.Clean()

	http.Handle("/", http.FileServer(http.Dir(*public)))
	http.Handle("/inbox", &server)

	httpServer := http.Server{
		Addr:           *bind,
		MaxHeaderBytes: *maxHeaderSize,
	}

	if *https {
		log.Fatal(httpServer.ListenAndServeTLS(*serverCert, *serverKey))
	} else {
		log.Fatal(httpServer.ListenAndServe())
	}
}
