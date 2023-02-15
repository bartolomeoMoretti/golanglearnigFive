package telegram

import (
	"errors"
	"golanglearningFive/clients/telegram"
	"golanglearningFive/events"
	"golanglearningFive/lib/e"
	"golanglearningFive/storage"
	"golanglearningFive/storage/files"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Meta struct {
	ChatID    int
	Username  string
	UserID    int64
	DateSent  int
	FirstName string
	LastName  string
}

func New(client *telegram.Client, storage files.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
		//добавить новый case если нужен другой метод, с котором раотает процессор
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username, meta.FirstName, meta.LastName, meta.UserID, meta.DateSent); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:    upd.Message.Chat.ID,
			Username:  upd.Message.From.Username,
			UserID:    upd.Message.From.UserID,
			FirstName: upd.Message.From.FirstName,
			LastName:  upd.Message.From.LastName,
			DateSent:  upd.Message.DateSent,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
