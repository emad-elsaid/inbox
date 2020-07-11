package inbox

import (
	"context"
	"errors"
	"time"
)

// Mailboxes holds all inboxes and when to timeout inboxes
type Mailboxes struct {
	inboxes       map[string]*inbox
	InboxCapacity int
	InboxTimeout  time.Duration
}

// New creates new empty Mailboxes structure with default timeouts
func New() *Mailboxes {
	return &Mailboxes{
		inboxes:       map[string]*inbox{},
		InboxCapacity: 100,
		InboxTimeout:  time.Minute,
	}
}

var (
	// ErrorIncorrectPassword is used whenever an operation tried to validate
	// inbox password and the password doesn't match the one in the inbox
	ErrorIncorrectPassword = errors.New("Incorrect password")
	// ErrorInboxNotFound is used when an operation tries to access an inbox but
	// the inbox doesn't exist
	ErrorInboxNotFound = errors.New("Inbox not found")
)

// Get the oldest message from `to` inbox, making sure the inbox password
// matches, it returns the message sender and the message, and an error if occurred
// This will also will restart the timeout for this inbox
func (m *Mailboxes) Get(to, password string, ctx context.Context) (from string, message []byte, err error) {
	inbox, ok := m.inboxes[to]
	if !ok {
		inbox = newInbox(password, m.InboxCapacity)
		m.inboxes[to] = inbox
	}

	if !inbox.CheckPassword(password) {
		err = ErrorIncorrectPassword
		return
	}

	from, message = inbox.Get(ctx)
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
		fromInbox = newInbox(password, m.InboxCapacity)
		m.inboxes[from] = fromInbox
	}

	if !fromInbox.CheckPassword(password) {
		return ErrorIncorrectPassword
	}

	return toInbox.Put(from, msg)
}

// Clean will delete timed out inboxes and messages
func (m *Mailboxes) Clean() {
	inboxDeadline := time.Now().Add(m.InboxTimeout * -1)
	for k, v := range m.inboxes {
		if v.blocking == 0 && v.lastAccessedAt.Before(inboxDeadline) {
			delete(m.inboxes, k)
		}
	}
}
