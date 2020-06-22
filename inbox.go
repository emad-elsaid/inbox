package inbox

import (
	"errors"
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

type Mailboxes struct {
	inboxes        map[string]*inbox
	InboxTimeout   time.Duration
	MessageTimeout time.Duration
}

func New() *Mailboxes {
	return &Mailboxes{
		inboxes:        map[string]*inbox{},
		InboxTimeout:   time.Minute,
		MessageTimeout: time.Minute,
	}
}

var (
	ErrorIncorrectPassword = errors.New("Incorrect password")
	ErrorInboxNotFound     = errors.New("Inbox not found")
)

func (m *Mailboxes) Get(to, password string) (string, []byte, error) {
	inbox, ok := m.inboxes[to]
	if !ok {
		inbox = newInbox(password)
		m.inboxes[to] = inbox
	}

	if !inbox.CheckPassword(password) {
		return "", nil, ErrorIncorrectPassword
	}

	from, message := inbox.Get()
	inbox.lastAccessedAt = time.Now()
	return from, message, nil
}

func (m *Mailboxes) Put(from, to, password string, msg []byte) error {
	toInbox, ok := m.inboxes[to]
	if !ok {
		return ErrorInboxNotFound
	}

	fromInbox, ok := m.inboxes[from]
	if !ok {
		fromInbox = newInbox(password)
		m.inboxes[from] = fromInbox
	}

	if !fromInbox.CheckPassword(password) {
		return ErrorIncorrectPassword
	}

	toInbox.Put(from, msg)
	return nil
}

func (m *Mailboxes) Clean() {
	inboxDeadline := time.Now().Add(m.InboxTimeout * -1)
	messageDeadline := time.Now().Add(m.MessageTimeout * -1)
	for k, v := range m.inboxes {
		if v.lastAccessedAt.Before(inboxDeadline) {
			delete(m.inboxes, k)
		} else {
			v.Clean(messageDeadline)
		}
	}
}
