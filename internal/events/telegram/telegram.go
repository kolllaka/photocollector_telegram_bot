package telegram

import (
	"fmt"
	"time"

	"github.com/kolllaka/telegram_bot/internal/clients/telegram"
	"github.com/kolllaka/telegram_bot/internal/events"
	"github.com/kolllaka/telegram_bot/internal/model"
	"github.com/kolllaka/telegram_bot/internal/storage"
	"github.com/kolllaka/telegram_bot/pkg/logging"
)

type processor struct {
	logger  *logging.Logger
	tg      telegram.Client
	offset  int
	storage storage.Storage
	locale  model.Commands
}

type Meta struct {
	ChatID int
	User   User
}

type User struct {
	ID       int
	Username string
}

func New(
	logger *logging.Logger,
	client telegram.Client,
	storage storage.Storage,
	locale model.Commands,
) *processor {
	return &processor{
		logger:  logger,
		tg:      client,
		storage: storage,
		locale:  locale,
	}
}

func (p *processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't get events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, p.event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *processor) Process(event events.Event) error {
	switch event.Type {
	case events.TextMessage:
		return p.processTextMessage(event)
	case events.PhotoMessage, events.FileMessage:
		return p.processPhotoMessage(event)
	default:
		return events.ErrUnknownEventType
	}
}

func (p *processor) processPhotoMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err := p.doPhotoCmd(event, meta); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	return nil
}

func (p *processor) processTextMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err := p.doCmd(event, meta); err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta: %w", events.ErrUnknownMetaType)
	}

	return res, nil
}

func (p *processor) event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
		Date: fetchTime(upd),
	}

	switch updType {
	case events.FileMessage:
		res.Url = upd.Message.Document.ID
	case events.PhotoMessage:
		res.Url = upd.Message.Photos[len(upd.Message.Photos)-1].ID
	}

	res.Meta = Meta{
		ChatID: upd.Message.Chat.ID,
		User: User{
			ID:       upd.Message.From.ID,
			Username: upd.Message.From.Username,
		},
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	if upd.Message.Photos != nil {
		return events.PhotoMessage
	}

	if upd.Message.Document != nil {
		return events.FileMessage
	}

	return events.TextMessage
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchTime(upd telegram.Update) time.Time {
	return time.Unix(int64(upd.Message.Date), 0)
}
