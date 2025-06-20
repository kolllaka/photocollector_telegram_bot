package telegram

import (
	"errors"
	"fmt"

	"github.com/kolllaka/telegram_bot/internal/clients/telegram"
	"github.com/kolllaka/telegram_bot/internal/model"
)

type storage struct {
	chatID   int
	tgClient telegram.Client
	locale   model.Infos
}

func New(
	chatID int,
	tgClient telegram.Client,
	locale model.Infos,
) *storage {
	return &storage{
		chatID:   chatID,
		tgClient: tgClient,
		locale:   locale,
	}
}

func (s *storage) Save(pm *model.PhotoMessage) error {
	msg := fmt.Sprintf(s.locale.TemplateMessageToChannel, pm.User.UserName, pm.ID)

	if err := s.tgClient.SendPhoto(s.chatID, pm.Url, msg); err != nil {
		if errors.Is(err, telegram.ErrPhotoNotSend) {
			s.tgClient.SendMessage(pm.ID, s.locale.WarnSendPhotoMessage)
		}

		return err
	}

	s.tgClient.SendMessage(pm.ID, s.locale.SuccessMessage)

	return nil
}
func (s *storage) Remove(pm *model.PhotoMessage) error {
	panic("not implement")
}
