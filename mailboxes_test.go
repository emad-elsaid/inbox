package inbox

import (
	"context"
	"testing"
	"time"
)

func TestMailboxes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Run("Mailboxes.Put", func(t *testing.T) {
		m := New()
		m.Get("Bob", "bob secret", ctx)

		err := m.Put("Alice", "Bob", "alice secret", []byte("message"))
		if err != nil {
			t.Errorf("got %s, expected no error", err)
		}

		err = m.Put("Alice", "Bob", "incorrect secret", []byte("message"))
		if err != ErrorIncorrectPassword {
			t.Errorf("got %s, expect %s", err, ErrorIncorrectPassword)
		}

		err = m.Put("Alice", "Fred", "alice secret", []byte("message"))
		if err != ErrorInboxNotFound {
			t.Errorf("Got %s, expected %s", err, ErrorInboxNotFound)
		}

		m.InboxCapacity = 0
		m.Get("BobFull", "bob secret", ctx)
		err = m.Put("Alice", "BobFull", "alice secret", []byte("message"))
		if err != ErrorInboxIsFull {
			t.Errorf("Got %s, expected %s", err, ErrorInboxIsFull)
		}

	})

	t.Run("Mailboxes.Get", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()

		m := New()
		from, msg, err := m.Get("Bob", "Bob secret", ctx)
		if from != "" {
			t.Errorf("Got %s, expected empty string", from)
		}

		if string(msg) != "" {
			t.Errorf("Got %s, expected empty string", msg)
		}

		m.Put("Alice", "Bob", "alice secret", []byte("hello"))
		from, msg, err = m.Get("Bob", "Bob secret", context.Background())
		if from != "Alice" {
			t.Errorf("Got %s, expected Alice", from)
		}

		if string(msg) != "hello" {
			t.Errorf("Got %s, expected hello", msg)
		}

		if err != nil {
			t.Errorf("Got %s, expected no error", err)
		}

		from, msg, err = m.Get("Bob", "wrong secret", context.Background())
		if from != "" {
			t.Errorf("Got %s, expected empty string", from)
		}

		if string(msg) != "" {
			t.Errorf("Got %s, expected empty string", msg)
		}

		if err != ErrorIncorrectPassword {
			t.Errorf("Got %s, expected %s", err, ErrorIncorrectPassword)
		}
	})

	t.Run("Mailboxes.Clean", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		m := New()
		m.InboxTimeout = 0
		m.Get("Alice", "secret", ctx)
		m.Get("Bob", "secret", ctx)
		m.Clean()
		err := m.Put("Bob", "Alice", "secret", []byte("hello"))
		m.Clean()
		if err != ErrorInboxNotFound {
			t.Errorf("Got %s, expected %s", err, ErrorInboxNotFound)
		}

		t.Run("When one of the inboxes are blocking a Get it shouldn't be deleted", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			m := New()
			m.InboxTimeout = 0
			go m.Get("Alice", "secret", ctx)
			time.Sleep(time.Millisecond)
			m.Clean()

			if len(m.inboxes) != 1 {
				t.Errorf("Inbox is deleted while Get is waiting: %d inboxes", len(m.inboxes))
			}
			cancel()
		})
	})
}
