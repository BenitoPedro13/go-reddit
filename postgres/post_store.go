package postgres

import (
	"fmt"

	goreddit "github.com/benitopedro13/go-reddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewPostStore(db *sqlx.DB) *PostStore {
	return &PostStore{DB: db}
}

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) Post(id uuid.UUID) (goreddit.Post, error) {
	var p goreddit.Post

	if err := s.DB.Get(&p, "SELECT * FROM posts WHERE id = $1", id); err != nil {
		return goreddit.Post{}, fmt.Errorf("failed to get post: %w", err)
	}

	return p, nil
}

func (s *PostStore) PostsByThread(threadID uuid.UUID) ([]goreddit.Post, error) {
	var pp []goreddit.Post

	if err := s.DB.Select(&pp, "SELECT * FROM posts WHERE thread_id = $1", threadID); err != nil {
		return []goreddit.Post{}, fmt.Errorf("failed to get posts: %w", err)
	}

	return pp, nil
}

func (s *PostStore) CreatePost(p *goreddit.Post) error {
	if err := s.DB.Get(p, "INSERT INTO posts VALUES ($1, $2, $3, $4, $5) RETURNING *",
		p.ID,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes); err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

func (s *PostStore) UpdatePost(p *goreddit.Post) error {
	if err := s.DB.Get(p, "UPDATE posts SET thread_id = $1, title = $2, content = $3, votes = $4 WHERE id = $5 RETURNING *",
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes,
		p.ID); err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := s.DB.Exec("DELETE FROM posts WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}
