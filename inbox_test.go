package inbox

import (
	"sync"
	"testing"
	"time"
)

func TestInbox(t *testing.T) {
	t.Run("Inbox.Get", func(t *testing.T) {
		i := newInbox("password")
		i.Put("Joe", []byte("hello"))
		if len(i.messages) != 1 {
			t.Errorf("len(messages) = %d; want 1", len(i.messages))
		}
	})

	t.Run("Inbox.Put", func(t *testing.T) {
		i := newInbox("password")

		from, msg := i.Get()
		if from != "" {
			t.Errorf("from = %s; want empty string", from)
		}

		if len(msg) != 0 {
			t.Errorf("message = %s; want empty bytes", msg)
		}

		i.Put("Joe", []byte("hello"))
		from, msg = i.Get()
		if from != "Joe" {
			t.Errorf("from = %s; want Joe", from)
		}

		if string(msg) != "hello" {
			t.Errorf("message = %s; want hello", msg)
		}
	})

	t.Run("Inbox.IsEmpty", func(t *testing.T) {
		i := newInbox("password")
		if !i.IsEmpty() {
			t.Errorf("expect inbox to be empty but it wasn't")
		}

		i.Put("Bob", []byte("message"))
		if i.IsEmpty() {
			t.Errorf("Expect inbox not to be empty but it was found empty")
		}
	})

	t.Run("Inbox.Clean", func(t *testing.T) {
		i := newInbox("password")
		i.Put("Joe", []byte("hello"))
		i.Clean(time.Now())

		from, msg := i.Get()
		if from != "" {
			t.Errorf("from = %s; want empty string", from)
		}

		if len(msg) != 0 {
			t.Errorf("message = %s; want empty bytes", msg)
		}

		i.Put("Joe", []byte("hello"))
		time.Sleep(time.Millisecond * 2)
		i.Put("Frank", []byte("world"))
		i.Clean(time.Now().Add(time.Millisecond * -1))

		from, msg = i.Get()
		if from != "Frank" {
			t.Errorf("from = %s; want Frank", from)
		}

		if string(msg) != "world" {
			t.Errorf("message = %s; want world", msg)
		}
	})
}

func BenchmarkInboxPut(b *testing.B) {
	i := newInbox("password")
	for n := 0; n < b.N; n++ {
		i.Put("Alice", []byte("Hello"))
	}
}

var (
	from string
	msg  []byte
)

func BenchmarkInboxPutThenGet(b *testing.B) {
	i := newInbox("password")
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for n := 0; n < b.N; n++ {
			i.Put("Alice", []byte("Hello"))
		}
		wg.Done()
	}()

	go func() {
		for n := 0; n < b.N; n++ {
			from, msg = i.Get()
		}
		wg.Done()
	}()

	wg.Wait()
}
