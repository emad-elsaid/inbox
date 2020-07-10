package inbox

import (
	"context"
	"testing"
)

func TestInbox(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Run("Inbox.Get", func(t *testing.T) {
		i := newInbox("password", 100)
		i.Put("Joe", []byte("hello"))
		if len(i.messages) != 1 {
			t.Errorf("len(messages) = %d; want 1", len(i.messages))
		}
	})

	t.Run("Inbox.Put", func(t *testing.T) {
		i := newInbox("password", 100)

		from, msg := i.Get(ctx)
		if from != "" {
			t.Errorf("from = %s; want empty string", from)
		}

		if len(msg) != 0 {
			t.Errorf("message = %s; want empty bytes", msg)
		}

		i.Put("Joe", []byte("hello"))
		from, msg = i.Get(context.Background())
		if from != "Joe" {
			t.Errorf("from = %s; want Joe", from)
		}

		if string(msg) != "hello" {
			t.Errorf("message = %s; want hello", msg)
		}
	})

	t.Run("Inbox.IsEmpty", func(t *testing.T) {
		i := newInbox("password", 100)
		if !i.IsEmpty() {
			t.Errorf("expect inbox to be empty but it wasn't")
		}

		i.Put("Bob", []byte("message"))
		if i.IsEmpty() {
			t.Errorf("Expect inbox not to be empty but it was found empty")
		}
	})
}

func BenchmarkInboxPut(b *testing.B) {
	i := newInbox("password", 100)
	for n := 0; n < b.N; n++ {
		i.Put("Alice", []byte("Hello"))
	}
}

var (
	from string
	msg  []byte
)

func BenchmarkInboxPutThenGet(b *testing.B) {
	i := newInbox("password", 100)
	for n := 0; n < b.N; n++ {
		i.Put("Alice", []byte("Hello"))
		from, msg = i.Get(context.Background())
	}
}
