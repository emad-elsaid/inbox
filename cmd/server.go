package main

import (
	"inbox"
	"log"
	"net/http"
	"time"
)

func main() {
	server := inbox.Server{
		Mailboxes:       inbox.New(),
		CleanupInterval: time.Second,
	}
	go server.Clean()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/inbox", server)
	log.Fatal(http.ListenAndServeTLS("0.0.0.0:3000", "server.crt", "server.key", nil))
}
