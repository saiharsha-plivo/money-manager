package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/saiharsha/money-manager/pkg/validator"
)

type CommentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

type Comment struct {
	ID          int64     `json:"id"`
	RecordID    int64     `json:"record_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Version     int64     `json:"version"`
}

var (
	ErrDuplicateComment = errors.New("duplicate comment")
)

func (c *CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (record_id, description)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, comment.RecordID, comment.Description).Scan(&comment.ID, &comment.CreatedAt, &comment.Version)
	if err != nil {
		c.ErrorLog.Print(err)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "comments_pkey"`:
			return ErrDuplicateComment
		default:
			return err
		}
	}

	return nil
}

func (c *CommentModel) Update(comment *Comment) error {
	query := `
		UPDATE comments
		SET description = $1, version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version
	`

	args := []interface{}{comment.Description, comment.ID, comment.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
	if err != nil {
		c.ErrorLog.Print(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (c *CommentModel) Delete(id int64) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		c.ErrorLog.Print(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.ErrorLog.Print(err)
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (c *CommentModel) GetByID(id int64) (*Comment, error) {
	query := `
		SELECT id, record_id, description, created_at, version
		FROM comments
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var comment Comment

	err := c.DB.QueryRowContext(ctx, query, id).Scan(&comment.ID, &comment.RecordID, &comment.Description, &comment.CreatedAt, &comment.Version)
	if err != nil {
		c.ErrorLog.Print(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &comment, nil
}

func (c *CommentModel) GetAll(recordID int64, filters Filters) ([]*Comment, error) {
	query := `
		SELECT id, record_id, description, created_at, version
		FROM comments
		WHERE record_id = $1
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3
	`

	args := []interface{}{recordID, filters.limit(), filters.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		c.ErrorLog.Print(err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	defer rows.Close()

	comments := []*Comment{}

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.RecordID, &comment.Description, &comment.CreatedAt, &comment.Version)
		if err != nil {
			c.ErrorLog.Print(err)
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		c.ErrorLog.Print(err)
		return nil, err
	}

	return comments, nil
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	v.Check(comment.Description != "", "description", "must be provided")
	v.Check(comment.RecordID != 0, "record_id", "must be provided")
}
