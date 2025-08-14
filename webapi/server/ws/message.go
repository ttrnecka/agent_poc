package ws

import "github.com/google/uuid"

type MessageType = int

// TYPE = 1 - client online announcement
// TYPE = 2 - client offline announcement
// TYPE = 10 - policy refresh signal

// collector messages
const (
	MSG_COLLECTOR_ONLINE  MessageType = iota + 1 // collector online announcement
	MSG_COLLECTOR_OFFLINE                        // collector offline announcement
)

// policy messages
const (
	MSG_POLICY_REFRESH MessageType = iota + 10 // policy refresh signal
)

// probe messages
const (
	MSG_PROBE_START        MessageType = iota + 20 // start probe signal
	MSG_PROBE_RUNNING                              // probe is running
	MSG_PROBE_FINISHED_OK                          // probe finished ok
	MSG_PROBE_FINISHED_ERR                         // probe finished with error
	MSG_PROBE_DATA
)

type Message struct {
	Type        MessageType
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
