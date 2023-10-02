package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"telegram-bot/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, username string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var (
	ErrNoSavedPage = errors.New("no saved page")
)

type Page struct {
	Url      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.Url); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	// чтобы правильно преобразовать массив байтов в строку
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
