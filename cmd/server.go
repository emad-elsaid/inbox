package main

import (
	"fmt"
	. "inbox"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	inboxes := map[string]*Inbox{}
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Inboxes: %d\n", len(inboxes))
	})

	http.HandleFunc("/inbox", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			inboxGet(inboxes, w, r)
		case http.MethodPost:
			inboxPost(inboxes, w, r)
		}
	})

	log.Fatal(http.ListenAndServeTLS("0.0.0.0:3000", "server.crt", "server.key", nil))
}

func inboxGet(inboxes map[string]*Inbox, w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	password := r.FormValue("password")
	inbox, ok := inboxes[to]
	if !ok {
		inbox = New(password)
		inboxes[to] = inbox
	}

	if !inbox.CheckPassword(password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	from, message := inbox.Get()
	w.Header().Add("X-From", from)
	fmt.Fprint(w, string(message))
}

func inboxPost(inboxes map[string]*Inbox, w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	password := r.FormValue("password")
	message, _ := ioutil.ReadAll(r.Body)

	toInbox, ok := inboxes[to]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	fromInbox, ok := inboxes[from]
	if !ok {
		fromInbox = New(password)
		inboxes[to] = fromInbox
	}

	if !fromInbox.CheckPassword(password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	toInbox.Put(from, message)
}
