package postgres

import (
	"fmt"

	goreddit "github.com/benitopedro13/go-reddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func NewCommentStore(db *sqlx.DB) *CommentStore {
	return &CommentStore{DB: db}
}

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var c goreddit.Comment

	if err := s.DB.Get(&c, "SELECT * FROM comments WHERE id = $1", id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("failed to get comment: %w", err)
	}

	return c, nil
}

func (s *CommentStore) CommentsByPost(postID uuid.UUID) ([]goreddit.Comment, error) {
	var cc []goreddit.Comment

	if err := s.DB.Select(&cc, "SELECT * FROM comments WHERE post_id = $1", postID); err != nil {
		return []goreddit.Comment{}, fmt.Errorf("failed to get comments: %w", err)
	}

	return cc, nil
}

func (s *CommentStore) CreateComment(c *goreddit.Comment) error {
	if err := s.DB.Get(c, "INSERT INTO comments VALUES ($1, $2, $3, $4) RETURNING *",
		c.ID,
		c.PostID,
		c.Content,
		c.Votes,
	); err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

func (s *CommentStore) UpdateComment(c *goreddit.Comment) error {
	if err := s.DB.Get(c, "UPDATE comments SET post_id = $1, content = $2, votes = $3 WHERE id = $4 RETURNING *",
		c.PostID,
		c.Content,
		c.Votes,
		c.ID); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := s.DB.Exec("DELETE FROM comments WHERE id = $1", id); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
