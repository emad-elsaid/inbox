package main

import (
	"fmt"
	"inbox"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	inboxes := inbox.New()
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", http.StripPrefix("/", fs))

	http.HandleFunc("/inbox", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			inboxGet(inboxes, w, r)
		case http.MethodPost:
			inboxPost(inboxes, w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	log.Fatal(http.ListenAndServeTLS("0.0.0.0:3000", "server.crt", "server.key", nil))
}

func inboxGet(mailboxes *inbox.Mailboxes, w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	password := r.FormValue("password")

	from, message, err := mailboxes.Get(to, password)
	if err!=nil {
		switch err {
		case inbox.ErrorIncorrectPassword:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Add("X-From", from)
	fmt.Fprint(w, string(message))
}

func inboxPost(mailboxes *inbox.Mailboxes, w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	password := r.FormValue("password")
	message, _ := ioutil.ReadAll(r.Body)

	if err := mailboxes.Put(from, to, password, message); err!=nil{
		switch err {
		case inbox.ErrorIncorrectPassword:
			w.WriteHeader(http.StatusUnauthorized)
		case inbox.ErrorInboxNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
