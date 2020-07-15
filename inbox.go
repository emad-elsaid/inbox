package inbox

import (
	"context"
	"errors"
	"sync/atomic"
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
	blocking       int32
	cancelContext  context.CancelFunc
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

func (i *inbox) Get(ctx *context.Context) (from string, msg []byte) {
	i.lastAccessedAt = time.Now()

	if ctx != nil {
		return i.getWithContext(ctx)
	} else {
		return i.getWithoutContext()
	}
}

func (i *inbox) getWithContext(ctx *context.Context) (from string, msg []byte) {
	atomic.AddInt32(&i.blocking, 1)

	wrapperCtx, cancel := context.WithCancel(*ctx)
	if i.cancelContext != nil {
		i.cancelContext()
	}
	i.cancelContext = cancel

	select {
	case message := <-i.messages:
		from = message.from
		msg = message.message
	case <-wrapperCtx.Done():
	}

	atomic.AddInt32(&i.blocking, -1)
	i.lastAccessedAt = time.Now()

	return
}

func (i *inbox) getWithoutContext() (from string, msg []byte) {
	select {
	case message := <-i.messages:
		return message.from, message.message
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

func (i *inbox) Locked() bool {
	return atomic.LoadInt32(&i.blocking) > 0
}
