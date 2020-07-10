package inbox

import (
	"errors"
	"time"
)

var (
	ErrorInboxIsFull = errors.New("Inbox is full")
)

type message struct {
	createdAt time.Time
	from      string
	message   []byte
}

type inbox struct {
	lastAccessedAt time.Time
	password       string
	messages       chan *message
}

func newInbox(password string, size int) *inbox {
	return &inbox{
		lastAccessedAt: time.Now(),
		password:       password,
		messages:       make(chan *message, size),
	}
}

func (i *inbox) Put(from string, msg []byte) error {
	select {
	case i.messages <- &message{
		createdAt: time.Now(),
		from:      from,
		message:   msg,
	}:
		return nil
	default:
		return ErrorInboxIsFull
	}
}

func (i *inbox) Get() (from string, msg []byte) {
	i.lastAccessedAt = time.Now()

	select {
	case msg := <-i.messages:
		return msg.from, msg.message
	default:
		return
	}
}

func (i *inbox) IsEmpty() bool {
	return len(i.messages) == 0
}

func (i *inbox) CheckPassword(password string) bool {
	return i.password == password
}
