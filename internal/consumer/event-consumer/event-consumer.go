package event_consumer

import (
	"time"

	"github.com/kolllaka/telegram_bot/internal/events"
	"github.com/kolllaka/telegram_bot/pkg/logging"
)

type consumer struct {
	logger    *logging.Logger
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

type Consumer interface {
	Start() error
}

func New(
	logger *logging.Logger,
	fetcher events.Fetcher,
	processor events.Processor,
	batchSize int,
) Consumer {
	return &consumer{
		logger:    logger,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.logger.Error("consumer", logging.ErrAttr(err))

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handlerEvents(gotEvents); err != nil {
			c.logger.Error("handlerEvents", logging.ErrAttr(err))

			continue
		}

	}
}

func (c *consumer) handlerEvents(events []events.Event) error {
	for _, event := range events {
		c.logger.Debug("got new event", logging.AnyAttr("event", event))

		if err := c.processor.Process(event); err != nil {
			c.logger.Warn("can't handler event", logging.ErrAttr(err))

			continue
		}
	}

	return nil
}
