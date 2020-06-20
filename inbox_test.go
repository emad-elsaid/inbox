package inbox

import "testing"

func TestPut(t *testing.T) {
	i := New("password")
	i.Put("Joe", []byte("hello"))
	if len(i.messages) != 1 {
		t.Errorf("len(messages) = %d; want 1", len(i.messages))
	}
}

func TestGet(t *testing.T) {
	i := New("password")

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
}
