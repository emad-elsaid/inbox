package inbox

import (
	"sync"
	"time"
)

type message struct {
	createdAt time.Time
	from      string
	message   []byte
}

type inbox struct {
	sync.RWMutex
	lastAccessedAt time.Time
	password       string
	messages       []message
}

func newInbox(password string) *inbox {
	return &inbox{
		lastAccessedAt: time.Now(),
		password:       password,
	}
}

func (i *inbox) Put(from string, msg []byte) {
	i.Lock()
	defer i.Unlock()

	i.messages = append(i.messages, message{
		createdAt: time.Now(),
		from:      from,
		message:   msg,
	})
}

func (i *inbox) Get() (from string, msg []byte) {
	i.RLock()
	defer i.RUnlock()

	i.lastAccessedAt = time.Now()

	if i.IsEmpty() {
		return
	}

	from = i.messages[0].from
	msg = i.messages[0].message
	i.messages = i.messages[1:]
	return
}

func (i *inbox) IsEmpty() bool {
	return len(i.messages) == 0
}

func (i *inbox) CheckPassword(password string) bool {
	return i.password == password
}

func (i *inbox) Clean(before time.Time) {
	i.Lock()
	defer i.Unlock()

	cutUntil := 0
	for ; cutUntil < len(i.messages) && i.messages[cutUntil].createdAt.Before(before); cutUntil++ {
	}
	i.messages = i.messages[cutUntil:]
}
