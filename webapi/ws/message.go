package ws

import "github.com/google/uuid"

type MessageType = int

// TYPE = 1 - client online announcement
// TYPE = 2 - client offline announcement
// TYPE = 10 - policy refresh signal

// collector messages
const (
	MSG_ONLINE  MessageType = iota + 1 // collector online announcement
	MSG_OFFLINE                        // collector offline announcement
)

// probe messages
const (
	MSG_REFRESH MessageType = iota + 10 // probe refresh signal
)

// session messages
const (
	MSG_RUN          MessageType = iota + 20 // run probe signal
	MSG_RUNNING                              // probe is running
	MSG_FINISHED_OK                          // probe finished ok
	MSG_FINISHED_ERR                         // probe finished with error
	MSG_DATA
)

type Message struct {
	Type        int
	Source      string
	Destination string
	Text        string
	Session     uuid.UUID `json:"Session,omitempty,omitzero"`
}

func NewMessage(t MessageType, source, destination, text string) Message {
	return Message{
		Type:        t,
		Source:      source,
		Destination: destination,
		Text:        text,
	}
}
