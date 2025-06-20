package telegram

import (
	"strings"
	"time"

	"github.com/kolllaka/telegram_bot/internal/events"
	"github.com/kolllaka/telegram_bot/internal/model"
	"github.com/kolllaka/telegram_bot/pkg/logging"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
)

func (p *processor) doCmd(event events.Event, meta Meta) error {
	text := strings.TrimSpace(event.Text)

	p.logger.Info(
		"got new text command",
		logging.StringAttr("text", text),
		logging.StringAttr("username", meta.User.Username),
		logging.IntAttr("username id", meta.User.ID),
	)

	switch text {
	case StartCmd:
		p.tg.SendMessage(meta.ChatID, p.locale.StartMessage)
	case HelpCmd:
		p.tg.SendMessage(meta.ChatID, p.locale.HelpMessage)
	default:
	}

	return nil
}

func (p *processor) doPhotoCmd(event events.Event, meta Meta) error {
	text := strings.TrimSpace(event.Text)

	p.logger.Info(
		"got new photo command",
		logging.StringAttr("text", text),
		logging.StringAttr("username", meta.User.Username),
		logging.IntAttr("username id", meta.User.ID),
	)

	photoMessage := model.PhotoMessage{
		ID: meta.ChatID,
		User: model.User{
			ID:       int64(meta.User.ID),
			UserName: meta.User.Username,
		},
		Url:  event.Url,
		Date: time.Now(),
	}

	if err := p.storage.Save(&photoMessage); err != nil {
		return err
	}

	return nil
}
