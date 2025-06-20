package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	tgBotHost = "api.telegram.org"

	FileMethod      = "file"
	GetFileMethod   = "getFile"
	sendPhotoMethod = "sendPhoto"
)

type Client interface {
	Updates(offset int, limit int) ([]Update, error)
	SendMessage(chatID int, text string) error
	SendPhoto(chatID int, fileID string, caption string) error
	FileLink(fileID string) string
}

type client struct {
	host     string
	basePath string
	client   http.Client
}

func New(token string) Client {
	return &client{
		host:     tgBotHost,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func (c *client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	if _, err := c.doRequest(sendMessageMethod, q); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	return nil
}

func (c *client) SendPhoto(chatID int, fileID string, caption string) error {
	type respInfo struct {
		Ok bool `json:"ok,omitempty"`
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("photo", fileID)
	q.Add("caption", caption)

	resp, err := c.doRequest(sendPhotoMethod, q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	var ri respInfo
	json.Unmarshal(resp, &ri)

	if !ri.Ok {
		return ErrPhotoNotSend
	}

	return nil
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *client) doRequest(method string, q url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Link returns a full path to the download URL for a File.
//
// It requires the Bot token to create the link.
func (c *client) FileLink(fileID string) string {
	filePath := c.filePath(fileID)

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(FileMethod, c.basePath, filePath),
	}

	return u.String()
}

// https://api.telegram.org/bot{token}/getFile
func (c *client) filePath(fileID string) string {
	type Result struct {
		FileId       string `json:"file_id,omitempty"`
		FileUniqueId string `json:"file_unique_id,omitempty"`
		FileSize     int    `json:"file_size,omitempty"`
		FilePath     string `json:"file_path,omitempty"`
	}
	type file struct {
		Ok     bool   `json:"ok,omitempty"`
		Result Result `json:"result,omitempty"`
	}

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, GetFileMethod),
	}

	data, _ := json.Marshal(map[string]string{
		"file_id": fileID,
	})

	req, _ := http.NewRequest("POST", u.String(), bytes.NewBuffer(data))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var fileInfo file
	json.Unmarshal(body, &fileInfo)

	return fileInfo.Result.FilePath
}
