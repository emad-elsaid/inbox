package inbox

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Server is an HTTP handler struct that holds all mailboxes in memory and when
// to clean them
type Server struct {
	Mailboxes       *Mailboxes
	CleanupInterval time.Duration
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.inboxGet(w, r)
	case http.MethodPost:
		s.inboxPost(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s Server) inboxGet(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	password := r.FormValue("password")

	from, message, err := s.Mailboxes.Get(to, password)
	if err != nil {
		switch err {
		case ErrorIncorrectPassword:
			w.WriteHeader(http.StatusUnauthorized)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if len(from) > 0 {
		w.Header().Add("X-From", from)
	}

	if len(message) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	fmt.Fprint(w, string(message))
}

func (s Server) inboxPost(w http.ResponseWriter, r *http.Request) {
	from := r.FormValue("from")
	to := r.FormValue("to")
	password := r.FormValue("password")
	message, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err := s.Mailboxes.Put(from, to, password, message); err != nil {
		switch err {
		case ErrorIncorrectPassword:
			w.WriteHeader(http.StatusUnauthorized)
		case ErrorInboxNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

// Clean will delete old inboxes and old messages periodically with a interval
// of CleanupInterval
func (s Server) Clean() {
	for {
		s.Mailboxes.Clean()
		time.Sleep(s.CleanupInterval)
	}
}
