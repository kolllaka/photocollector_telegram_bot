package model

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

type PhotoMessage struct {
	ID   int
	User User
	Url  string
	Date time.Time
}

type User struct {
	ID       int64
	UserName string
}

func (pm PhotoMessage) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, pm.Url); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	if _, err := io.WriteString(h, pm.User.UserName); err != nil {
		return "", fmt.Errorf("can't calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
