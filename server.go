package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	MESSAGES_SENDER   = [][]byte{}
	MESSAGES_RECEIVER = [][]byte{}
)

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		MESSAGES_SENDER = [][]byte{}
		MESSAGES_RECEIVER = [][]byte{}

		view, _ := ioutil.ReadFile("views/send.html")
		fmt.Fprint(w, string(view))
	})

	http.HandleFunc("/receive", func(w http.ResponseWriter, r *http.Request) {
		view, _ := ioutil.ReadFile("views/receive.html")
		fmt.Fprint(w, string(view))
	})

	http.HandleFunc("/from/receiver", func(w http.ResponseWriter, r *http.Request) {
		message, _ := ioutil.ReadAll(r.Body)
		MESSAGES_SENDER = append(MESSAGES_SENDER, message)
	})

	http.HandleFunc("/inbox/receiver", func(w http.ResponseWriter, r *http.Request) {
		if len(MESSAGES_RECEIVER) == 0 {
			return
		}

		fmt.Fprint(w, string(MESSAGES_RECEIVER[0]))
		MESSAGES_RECEIVER = MESSAGES_RECEIVER[1:]
	})

	http.HandleFunc("/from/sender", func(w http.ResponseWriter, r *http.Request) {
		message, _ := ioutil.ReadAll(r.Body)
		MESSAGES_RECEIVER = append(MESSAGES_RECEIVER, message)
	})

	http.HandleFunc("/inbox/sender", func(w http.ResponseWriter, r *http.Request) {
		if len(MESSAGES_SENDER) == 0 {
			return
		}

		fmt.Fprint(w, string(MESSAGES_SENDER[0]))
		MESSAGES_SENDER = MESSAGES_SENDER[1:]
	})

	log.Fatal(http.ListenAndServeTLS("0.0.0.0:3000", "server.crt", "server.key", nil))
}
