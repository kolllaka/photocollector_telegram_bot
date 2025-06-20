package events

import "errors"

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Type int

const (
	Unknown Type = iota
	TextMessage
	PhotoMessage
	FileMessage
)

type Event struct {
	Type Type
	Url  string
	Text string
	Meta any
}
