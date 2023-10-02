package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"telegram-bot/clients/telegram"
	"telegram-bot/lib/e"
	"telegram-bot/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatId int, username string) (err error) {
	defer func() { err = e.Wrap("can't doCmd", err) }()

	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(ctx, chatId, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatId, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatId)
	case StartCmd:
		return p.sendHello(ctx, chatId)
	default:
		return p.tg.SendMessage(chatId, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageUrl string, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: save page", err) }()

	sendMsg := newMessageSender(chatID, p.tg)

	page := &storage.Page{
		Url:      pageUrl,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: send random", err) }()

	sendMsg := newMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPage) {
		return sendMsg(msgNoSavedPages)
	}
	if err := sendMsg(page.Url); err != nil {
		return err
	}
	return p.storage.Remove(ctx, page)
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func newMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
