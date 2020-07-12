package inbox

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// Server is an HTTP handler struct that holds all mailboxes in memory and when
// to clean them
type Server struct {
	CORS            bool
	Mailboxes       *Mailboxes
	CleanupInterval time.Duration
	MaxBodySize     int64
	LongPolling     bool
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.MaxBodySize)

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.inboxGet(w, r)
	case http.MethodPost:
		s.inboxPost(w, r)
	case http.MethodOptions:
		s.inboxOptions(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) inboxGet(w http.ResponseWriter, r *http.Request) {
	s.writeCORS(w)

	to, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestCtx := r.Context()
	var ctx *context.Context

	if s.LongPolling {
		ctx = &requestCtx
	}

	from, message, err := s.Mailboxes.Get(to, password, ctx)
	if err != nil {
		switch err {
		case ErrorIncorrectPassword:
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
		}
		return
	}

	w.Header().Add("X-From", from)
	w.Write(message)
}

func (s *Server) inboxPost(w http.ResponseWriter, r *http.Request) {
	s.writeCORS(w)

	from, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	message, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	r.Body.Close()

	to := r.FormValue("to")
	if err := s.Mailboxes.Put(from, to, password, message); err != nil {
		switch err {
		case ErrorIncorrectPassword:
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
		case ErrorInboxNotFound:
			w.WriteHeader(http.StatusNotFound)
		case ErrorInboxIsFull:
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

func (s *Server) inboxOptions(w http.ResponseWriter, r *http.Request) {
	s.writeCORS(w)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) writeCORS(w http.ResponseWriter) {
	if !s.CORS {
		return
	}

	headers := w.Header()
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	headers.Set("Access-Control-Allow-Origin", "*")
	headers.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	headers.Set("Access-Control-Allow-Credentials", "true")
	headers.Set("Access-Control-Allow-Headers", "Authorization")
	headers.Set("Access-Control-Expose-Headers", "X-From")
}

// Clean will delete old inboxes periodically with an interval of CleanupInterval
func (s Server) Clean() {
	for {
		s.Mailboxes.Clean()
		time.Sleep(s.CleanupInterval)
	}
}
