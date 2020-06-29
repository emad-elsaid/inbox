package inbox

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	handler := Server{Mailboxes: New()}

	t.Run("GET", func(t *testing.T) {
		t.Run("wrong body format", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/inbox", nil)
			req.URL.RawQuery = "someWeird%<asdfa\nstuf*))\n-f"
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}

			if rr.Body.String() != "" {
				t.Errorf("handler returned unexpected body: got %v wanted empty string", rr.Body.String())
			}
		})

		t.Run("First request", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/inbox", nil)
			req.URL.RawQuery = url.Values{"to": {"Alice"}, "password": {"secret"}}.Encode()
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNoContent {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
			}

			if rr.Body.String() != "" {
				t.Errorf("handler returned unexpected body: got %v wanted empty string", rr.Body.String())
			}
		})

		t.Run("When password is incorrect", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/inbox", nil)
			req.URL.RawQuery = url.Values{"to": {"Alice"}, "password": {"incorrect secret"}}.Encode()
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
			}

			if rr.Body.String() != "" {
				t.Errorf("handler returned unexpected body: got %v wanted empty string", rr.Body.String())
			}
		})

		t.Run("When a message is in the inbox", func(t *testing.T) {
			handler.Mailboxes.Put("Bob", "Alice", "secret", []byte("message"))

			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/inbox", nil)
			req.URL.RawQuery = url.Values{"to": {"Alice"}, "password": {"secret"}}.Encode()
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}


			if from := rr.HeaderMap.Get("X-From"); from != "Bob" {
				t.Errorf("handler returned wrong X-From header: got %v want %v", from, "Bob")
			}

			if rr.Body.String() != "message" {
				t.Errorf("handler returned unexpected body: got %v wanted %s", rr.Body.String(), "message")
			}
		})
	})

	t.Run("POST", func(t *testing.T) {
		t.Run("When password is correct and inbox exists", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/inbox", strings.NewReader("message"))
			req.URL.RawQuery = url.Values{"from": {"Bob"}, "to": {"Alice"}, "password": {"secret"}}.Encode()
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		})

		t.Run("When sending message to non-existent inbox", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/inbox", strings.NewReader("message"))
			req.URL.RawQuery = url.Values{"from": {"Bob"}, "to": {"Caty"}, "password": {"secret"}}.Encode()
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusNotFound {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
			}
		})

		t.Run("When password is incorrect", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/inbox", strings.NewReader("message"))
			req.URL.RawQuery = url.Values{"from": {"Bob"}, "to": {"Alice"}, "password": {"incorrect secret"}}.Encode()
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
			}
		})
	})
}
