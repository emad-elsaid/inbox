package inbox

import "time"

type message struct {
	createdAt time.Time
	from      string
	message   []byte
}

type Inbox struct {
	createdAt time.Time
	password  string
	messages  []message
}

func New(password string) *Inbox {
	return &Inbox{
		createdAt: time.Now(),
		password:  password,
	}
}

func (i *Inbox) Put(from string, msg []byte) {
	i.messages = append(i.messages, message{
		createdAt: time.Now(),
		from:      from,
		message:   msg,
	})
}

func (i *Inbox) Get() (string, []byte) {
	if len(i.messages) == 0 {
		return "", []byte{}
	}

	message := i.messages[0]
	i.messages = i.messages[1:]
	return message.from, message.message
}

func (i *Inbox) CheckPassword(password string) bool {
	return i.password == password
}
