package postgres

import (
	"fmt"

	goreddit "github.com/benitopedro13/go-reddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewThreadStore(db *sqlx.DB) *ThreadStore {
	return &ThreadStore{DB: db}
}

type ThreadStore struct {
	*sqlx.DB
}

func (s *ThreadStore) Thread(id uuid.UUID) (goreddit.Thread, error) {
	var t goreddit.Thread

	if err := s.DB.Get(&t, "SELECT * FROM threads WHERE id = $1", id); err != nil {
		return goreddit.Thread{}, fmt.Errorf("failed to get thread: %w", err)
	}

	return t, nil
}

func (s *ThreadStore) Threads() ([]goreddit.Thread, error) {
	tt := []goreddit.Thread{}

	if err := s.DB.Select(&tt, "SELECT * FROM threads"); err != nil {
		return []goreddit.Thread{}, fmt.Errorf("failed to get threads: %w", err)
	}

	return tt, nil
}

func (s *ThreadStore) CreateThread(t *goreddit.Thread) error {
	if err := s.DB.Get(t, "INSERT INTO threads ($1, $2, $3) RETURNING *", t.ID, t.Title, t.Description); err != nil {
		return fmt.Errorf("failed to create thread: %w", err)
	}

	return nil
}

func (s *ThreadStore) UpdateThread(t *goreddit.Thread) error {
	if err := s.DB.Get(t, "UPDATE threads SET title = $1, description = $2 WHERE id = $3 RETURNING *", t.Title, t.Description, t.ID); err != nil {
		return fmt.Errorf("failed to update thread: %w", err)
	}

	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := s.DB.Exec("DELETE FROM threads WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete thread: %w", err)
	}

	return nil
}
