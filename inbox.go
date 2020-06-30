package inbox

import (
	"time"
)

type message struct {
	createdAt time.Time
	from      string
	message   []byte
}

type inbox struct {
	createdAt      time.Time
	lastAccessedAt time.Time
	password       string
	messages       []message
}

func newInbox(password string) *inbox {
	return &inbox{
		createdAt:      time.Now(),
		lastAccessedAt: time.Now(),
		password:       password,
	}
}

func (i *inbox) Put(from string, msg []byte) {
	i.messages = append(i.messages, message{
		createdAt: time.Now(),
		from:      from,
		message:   msg,
	})
}

func (i *inbox) Get() (string, []byte) {
	if len(i.messages) == 0 {
		return "", []byte{}
	}

	message := i.messages[0]
	i.messages = i.messages[1:]
	return message.from, message.message
}

func (i *inbox) CheckPassword(password string) bool {
	return i.password == password
}

func (i *inbox) Clean(before time.Time) {
	cutUntil := 0
	for ; cutUntil < len(i.messages) && i.messages[cutUntil].createdAt.Before(before); cutUntil++ {
	}
	i.messages = i.messages[cutUntil:]
}
