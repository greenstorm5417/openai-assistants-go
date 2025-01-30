package streaming

type EventType string

const (
	EventTypeMessage EventType = "message"
	EventTypeError   EventType = "error"
	EventTypeDone    EventType = "done"
)

type Event struct {
	Type    EventType   `json:"type"`
	Data    string      `json:"data"`
	Error   *ErrorData  `json:"error,omitempty"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StreamReader interface {
	Next() (*Event, error)
	Close() error
}