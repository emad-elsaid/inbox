package inbox

import (
	"time"
	"errors"
)
// Mailboxes holds all inboxes and when to timeout inboxes and messages
type Mailboxes struct {
	inboxes        map[string]*inbox
	InboxTimeout   time.Duration
	MessageTimeout time.Duration
}

// New creates new empty Mailboxes structure with default timeouts
func New() *Mailboxes {
	return &Mailboxes{
		inboxes:        map[string]*inbox{},
		InboxTimeout:   time.Minute,
		MessageTimeout: time.Minute,
	}
}

var (
	// ErrorIncorrectPassword is used whenever an operation tried to validate
	// inbox password and the password doesn't match the one in the inbox
	ErrorIncorrectPassword = errors.New("Incorrect password")
	// ErrorInboxNotFound is used when an operation tries to access an inbox but
	// the inbox doesn't exist
	ErrorInboxNotFound = errors.New("Inbox not found")
	// ErrorInboxIsEmpty is used if tried to get a message from empty inbox
	ErrorInboxIsEmpty = errors.New("Inbox is empty")
)

// Get the oldest message from `to` inbox, making sure the inbox password
// matches, it returns the message sender and the message, and an error if occurred
// This will also will restart the timeout for this inbox
func (m *Mailboxes) Get(to, password string) (from string, message []byte, err error) {
	inbox, ok := m.inboxes[to]
	if !ok {
		inbox = newInbox(password)
		m.inboxes[to] = inbox
	}

	if !inbox.CheckPassword(password) {
		err = ErrorIncorrectPassword
		return
	}

	if inbox.IsEmpty() {
		err = ErrorInboxIsEmpty
		return
	}

	from, message = inbox.Get()
	inbox.lastAccessedAt = time.Now()
	return
}

// Put will put a message `msg` at the end of `to` inbox from inbox owned by
// `from` if the password `password` matches the one stored in `from` inbox.
// If `from` Inbox doesn't exist it will be created with `password`.
// If `to` inbox doesn't exist it will return ErrorInboxNotFound
// When the `from` inbox exist and the password doesn't match it will return ErrorIncorrectPassword
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

// Clean will delete timed out inboxes and messages
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
