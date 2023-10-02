package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"telegram-bot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, err
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := "INSERT INTO pages(url, user_name) VALUES (?, ?);"

	if _, err := s.db.ExecContext(ctx, q, p.Url, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) PickRandom(ctx context.Context, username string) (*storage.Page, error) {
	q := "SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1;"

	var url string

	err := s.db.QueryRowContext(ctx, q, username).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrNoSavedPage
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		Url:      url,
		UserName: username,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := "DELETE FROM pages WHERE user_name = ? AND url = ?;"

	if _, err := s.db.ExecContext(ctx, q, p.UserName, p.Url); err != nil {
		return fmt.Errorf("can't remove pages %w", err)
	}

	return nil
}

func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	q := "SELECT (true) FROM pages WHERE user_name = ? AND url = ?;"

	var exists bool

	err := s.db.QueryRowContext(ctx, q, p.UserName, p.Url).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return true, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := "CREATE TABLE IF NOT EXISTS pages(url TEXT, user_name TEXT);"

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
