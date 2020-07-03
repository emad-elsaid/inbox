package inbox

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	handler := Server{CORS: true, Mailboxes: New()}

	t.Run("GET", func(t *testing.T) {
		t.Run("without authorization", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/inbox", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
			}

			if header := rr.HeaderMap.Get("WWW-Authenticate"); header != "Basic" {
				t.Errorf("handler returned unexpected WWW-Authenticate header: got %v wanted Basic", header)
			}

			if rr.Body.String() != "" {
				t.Errorf("handler returned unexpected body: got %v wanted empty string", rr.Body.String())
			}
		})

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
			req.SetBasicAuth("Alice", "secret")
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
			req.SetBasicAuth("Alice", "incorrect secret")
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
			req.SetBasicAuth("Alice", "secret")
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
		t.Run("without authorization", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/inbox", strings.NewReader("hello"))
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
			}

			if header := rr.HeaderMap.Get("WWW-Authenticate"); header != "Basic" {
				t.Errorf("handler returned unexpected WWW-Authenticate header: got %v wanted Basic", header)
			}

			if rr.Body.String() != "" {
				t.Errorf("handler returned unexpected body: got %v wanted empty string", rr.Body.String())
			}
		})

		t.Run("When password is correct and inbox exists", func(t *testing.T) {
			rr := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/inbox", strings.NewReader("message"))
			req.URL.RawQuery = url.Values{"to": {"Alice"}}.Encode()
			req.SetBasicAuth("Bob", "secret")
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
			req.URL.RawQuery = url.Values{"to": {"Caty"}}.Encode()
			req.SetBasicAuth("Bob", "secret")
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
			req.URL.RawQuery = url.Values{"to": {"Alice"}}.Encode()
			req.SetBasicAuth("Bob", "incorrect secret")
			if err != nil {
				t.Fatal(err)
			}

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusUnauthorized {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
			}
		})
	})

	t.Run("HEAD", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("HEAD", "/inbox", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}

		if rr.Body.String() != "" {
			t.Errorf("handler returned unexpected body: got %v wanted empty string", rr.Body.String())
		}
	})
}

// BENCHMARKS
var rr *httptest.ResponseRecorder

func newGetRequest(username, password string) *http.Request {
	req, err := http.NewRequest("GET", "/inbox", nil)
	req.SetBasicAuth(username, password)
	if err != nil {
		panic("error creating request")
	}
	return req
}

func BenchmarkServerGet(b *testing.B) {
	handler := Server{Mailboxes: New()}
	requests := []*http.Request{
		newGetRequest("Alice", "alicepassword"),
		newGetRequest("Bob", "bobpassword"),
		newGetRequest("Carole", "carolepassword"),
		newGetRequest("Dave", "davepassword"),
	}

	for n := 0; n < b.N; n++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, requests[n%len(requests)])
	}
}

func newPostRequest(username, password string, message []byte) *http.Request {
	req, err := http.NewRequest("POST", "/inbox", bytes.NewReader(message))
	req.URL.RawQuery = url.Values{"to": {"Alice"}}.Encode()
	req.SetBasicAuth(username, password)
	if err != nil {
		panic("error creating request")
	}
	return req
}

func BenchmarkServerPost(b *testing.B) {
	handler := Server{Mailboxes: New()}
	handler.Mailboxes.Get("Alice", "alicepassword")

	requests := []*http.Request{
		newPostRequest("Bob", "bobpassword", []byte("hello world bob")),
		newPostRequest("Carole", "carolepassword", []byte("hello world carole")),
		newPostRequest("Dave", "davepassword", []byte("hello world dave")),
	}

	for n := 0; n < b.N; n++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, requests[n%len(requests)])
	}
}
