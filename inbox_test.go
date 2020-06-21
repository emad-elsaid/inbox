package inbox

import "testing"

func TestInbox(t *testing.T) {
	t.Run("Inbox.Get", func(t *testing.T) {
		i := NewInbox("password")
		i.Put("Joe", []byte("hello"))
		if len(i.messages) != 1 {
			t.Errorf("len(messages) = %d; want 1", len(i.messages))
		}
	})

 t.Run("Inbox.Put", func(t *testing.T) {
	 i := NewInbox("password")

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
}

func TestMailboxes(t *testing.T) {
	t.Run("Mailboxes.Put", func(t *testing.T){
		m := New()
		m.Get("Bob", "bob secret")

		err := m.Put("Alice", "Bob", "alice secret", []byte("message"))
		if err!=nil{
			t.Errorf("got %s, expected no error", err)
		}

		err = m.Put("Alice", "Bob", "incorrect secret", []byte("message"))
		if err!=ErrorIncorrectPassword{
			t.Errorf("got %s, expect %s", err, ErrorIncorrectPassword)
		}

		err = m.Put("Alice", "Fred", "alice secret", []byte("message"))
		if err!=ErrorInboxNotFound {
			t.Errorf("Got %s, expected %s", err, ErrorInboxNotFound)
		}
	})

	t.Run("Mailboxes.Get", func(t *testing.T) {
		m := New()
		from, msg, err := m.Get("Bob", "Bob secret")
		if from != "" {
			t.Errorf("Got %s, expected empty string", from)
		}

		if string(msg) != "" {
			t.Errorf("Got %s, expected empty string", msg)
		}

		if err!= nil {
			t.Errorf("Got %s, expected no error", err)
		}

		m.Put("Alice", "Bob", "alice secret", []byte("hello"))
		from, msg, err = m.Get("Bob", "Bob secret")
		if from != "Alice" {
			t.Errorf("Got %s, expected Alice", from)
		}

		if string(msg) != "hello" {
			t.Errorf("Got %s, expected hello", msg)
		}

		if err!= nil {
			t.Errorf("Got %s, expected no error", err)
		}

		from, msg, err = m.Get("Bob", "wrong secret")
		if from != "" {
			t.Errorf("Got %s, expected empty string", from)
		}

		if string(msg) != "" {
			t.Errorf("Got %s, expected empty string", msg)
		}

		if err!= ErrorIncorrectPassword {
			t.Errorf("Got %s, expected %s", err, ErrorIncorrectPassword)
		}
	})
}

func BenchmarkInboxPut(b *testing.B) {
	i := NewInbox("password")
	for n := 0; n < b.N; n++ {
		i.Put("Alice", []byte("Hello"))
	}
}

var (
	from string
	msg []byte
)
func BenchmarkInboxPutThenGet(b *testing.B) {
	i := NewInbox("password")
	for n := 0; n < b.N; n++ {
		i.Put("Alice", []byte("Hello"))
		from, msg = i.Get()
	}
}
