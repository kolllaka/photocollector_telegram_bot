package storage

import (
	"github.com/kolllaka/telegram_bot/internal/model"
)

type Storage interface {
	Save(pm *model.PhotoMessage) error
	Remove(pm *model.PhotoMessage) error
}
