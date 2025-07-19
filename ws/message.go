package ws

type MessageType int

// TYPE = 1 - client online announcement
// TYPE = 2 - client offline announcement
// TYPE = 10 - policy refresh signal

// collector messages
const (
	MSG_ONLINE  MessageType = iota + 1 // collector online announcement
	MSG_OFFLINE                        // collector offline announcement
)

// policy messages
const (
	MSG_REFRESH MessageType = iota + 10 // Policy refresh signal
)

type Message struct {
	Type   int
	Source string
	Text   string
}

func NewMessage(t MessageType, source, text string) Message {
	return Message{
		Type:   int(t),
		Source: source,
		Text:   text,
	}
}
