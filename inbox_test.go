package inbox

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInbox(t *testing.T) {
	t.Run("Inbox.Get", func(t *testing.T) {
		i := newInbox("password", 100)
		i.Put("Joe", []byte("hello"))
		if len(i.messages) != 1 {
			t.Errorf("len(messages) = %d; want 1", len(i.messages))
		}
	})

	t.Run("Inbox.Put", func(t *testing.T) {
		i := newInbox("password", 100)

		from, msg := i.Get(nil)
		if from != "" {
			t.Errorf("from = %s; want empty string", from)
		}

		if len(msg) != 0 {
			t.Errorf("message = %s; want empty bytes", msg)
		}

		i.Put("Joe", []byte("hello"))
		from, msg = i.Get(nil)
		if from != "Joe" {
			t.Errorf("from = %s; want Joe", from)
		}

		if string(msg) != "hello" {
			t.Errorf("message = %s; want hello", msg)
		}

		t.Run("Waits for context", func(t *testing.T) {
			i := newInbox("password", 100)
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				from, msg := i.Get(&ctx)
				if from != "" {
					t.Errorf("from = %s; want empty string", from)
				}

				if len(msg) != 0 {
					t.Errorf("message = %s; want empty bytes", msg)
				}
			}()

			time.Sleep(time.Millisecond)
			cancel()
		})

		t.Run("Waits for a message", func(t *testing.T) {
			i := newInbox("password", 100)
			ctx, _ := context.WithCancel(context.Background())
			go func() {
				from, msg := i.Get(&ctx)
				if from != "Bob" {
					t.Errorf("from = %s; want Bob", from)
				}

				if string(msg) != "message" {
					t.Errorf("message = %s; want message", msg)
				}
			}()

			time.Sleep(time.Millisecond)
			i.Put("Bob", []byte("message"))
		})
	})

	t.Run("When two gets are waiting the newest gets themessage", func(t *testing.T) {
		i := newInbox("password", 100)
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			ctx, _ := context.WithCancel(context.Background())
			from, msg := i.Get(&ctx)
			if from != "" {
				t.Errorf("from = %s; want empty string", from)
			}

			if string(msg) != "" {
				t.Errorf("message = %s; want empty message", msg)
			}
			wg.Done()
		}()

		time.Sleep(time.Millisecond)

		go func() {
			ctx, _ := context.WithCancel(context.Background())
			from, msg := i.Get(&ctx)
			if from != "Bob" {
				t.Errorf("from = %s; want Bob", from)
			}

			if string(msg) != "message" {
				t.Errorf("message = %s; want message", msg)
			}
			wg.Done()
		}()

		time.Sleep(time.Millisecond)

		i.Put("Bob", []byte("message"))
		wg.Wait()
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
		from, msg = i.Get(nil)
	}
}
