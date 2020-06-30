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
	to, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	from, message, err := s.Mailboxes.Get(to, password)
	if err != nil {
		switch err {
		case ErrorIncorrectPassword:
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
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
	from, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	to := r.FormValue("to")
	message, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if err := s.Mailboxes.Put(from, to, password, message); err != nil {
		switch err {
		case ErrorIncorrectPassword:
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
		case ErrorInboxNotFound:
			w.WriteHeader(http.StatusNotFound)
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
