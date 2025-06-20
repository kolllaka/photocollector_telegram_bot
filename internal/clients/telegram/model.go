package telegram

import "errors"

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

var (
	ErrPhotoNotSend = errors.New("photo not send")
)

type UpdatesResponse struct {
	Ok     bool     `json:"ok,omitempty"`
	Result []Update `json:"result,omitempty"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message,omitempty"`
}

type IncomingMessage struct {
	Text     string    `json:"text,omitempty"`
	From     User      `json:"from,omitempty"`
	Chat     Chat      `json:"chat,omitempty"`
	Photos   []Photo   `json:"photo,omitempty"`
	Document *Document `json:"document,omitempty"`
	Date     int       `json:"date,omitempty"`
}

type Chat struct {
	ID int `json:"id,omitempty"`
}

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

type Photo struct {
	ID string `json:"file_id,omitempty"`
}

type Document struct {
	ID       string `json:"file_id,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}
