// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Book struct {
	ID             int32          `json:"id"`
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	Content        []byte         `json:"content"`
	CoverImage     []byte         `json:"cover_image"`
	CoverImageType sql.NullString `json:"cover_image_type"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
}