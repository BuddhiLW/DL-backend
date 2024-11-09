// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: books.sql

package db

import (
	"context"
	"database/sql"
)

const countBooks = `-- name: CountBooks :one
SELECT COUNT(*) AS count FROM books
`

func (q *Queries) CountBooks(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.countBooksStmt, countBooks)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createBook = `-- name: CreateBook :exec
INSERT INTO books (
    title,
    author,
    category,
    content,
    cover_image,
    cover_image_type
) VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateBookParams struct {
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	Content        []byte         `json:"content"`
	CoverImage     []byte         `json:"cover_image"`
	CoverImageType sql.NullString `json:"cover_image_type"`
}

func (q *Queries) CreateBook(ctx context.Context, arg CreateBookParams) error {
	_, err := q.exec(ctx, q.createBookStmt, createBook,
		arg.Title,
		arg.Author,
		arg.Category,
		arg.Content,
		arg.CoverImage,
		arg.CoverImageType,
	)
	return err
}

const deleteBook = `-- name: DeleteBook :exec
DELETE FROM books WHERE id = $1
`

func (q *Queries) DeleteBook(ctx context.Context, id int32) error {
	_, err := q.exec(ctx, q.deleteBookStmt, deleteBook, id)
	return err
}

const getBook = `-- name: GetBook :one
SELECT
    id,
    title,
    author,
    category,
    content,
    cover_image,
    cover_image_type
FROM books
WHERE id = $1
`

type GetBookRow struct {
	ID             int32          `json:"id"`
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	Content        []byte         `json:"content"`
	CoverImage     []byte         `json:"cover_image"`
	CoverImageType sql.NullString `json:"cover_image_type"`
}

func (q *Queries) GetBook(ctx context.Context, id int32) (GetBookRow, error) {
	row := q.queryRow(ctx, q.getBookStmt, getBook, id)
	var i GetBookRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Category,
		&i.Content,
		&i.CoverImage,
		&i.CoverImageType,
	)
	return i, err
}

const listBooks = `-- name: ListBooks :many
SELECT
    id,
    title,
    author,
    category,
    cover_image_type
FROM books
ORDER BY id
`

type ListBooksRow struct {
	ID             int32          `json:"id"`
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	CoverImageType sql.NullString `json:"cover_image_type"`
}

func (q *Queries) ListBooks(ctx context.Context) ([]ListBooksRow, error) {
	rows, err := q.query(ctx, q.listBooksStmt, listBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBooksRow
	for rows.Next() {
		var i ListBooksRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Category,
			&i.CoverImageType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchBooksByAuthor = `-- name: SearchBooksByAuthor :many
SELECT
    id,
    title,
    author,
    category,
    cover_image_type
FROM books
WHERE author LIKE '%' || $1 || '%'
ORDER BY author
`

type SearchBooksByAuthorRow struct {
	ID             int32          `json:"id"`
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	CoverImageType sql.NullString `json:"cover_image_type"`
}

func (q *Queries) SearchBooksByAuthor(ctx context.Context, dollar_1 sql.NullString) ([]SearchBooksByAuthorRow, error) {
	rows, err := q.query(ctx, q.searchBooksByAuthorStmt, searchBooksByAuthor, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchBooksByAuthorRow
	for rows.Next() {
		var i SearchBooksByAuthorRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Category,
			&i.CoverImageType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchBooksByTitle = `-- name: SearchBooksByTitle :many
SELECT
    id,
    title,
    author,
    category,
    cover_image_type
FROM books
WHERE title LIKE '%' || $1 || '%'
ORDER BY title
`

type SearchBooksByTitleRow struct {
	ID             int32          `json:"id"`
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	CoverImageType sql.NullString `json:"cover_image_type"`
}

func (q *Queries) SearchBooksByTitle(ctx context.Context, dollar_1 sql.NullString) ([]SearchBooksByTitleRow, error) {
	rows, err := q.query(ctx, q.searchBooksByTitleStmt, searchBooksByTitle, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchBooksByTitleRow
	for rows.Next() {
		var i SearchBooksByTitleRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Category,
			&i.CoverImageType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBook = `-- name: UpdateBook :exec
UPDATE books
SET
    title = $1,
    author = $2,
    category = $3,
    content = $4,
    cover_image = $5,
    cover_image_type = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $7
`

type UpdateBookParams struct {
	Title          string         `json:"title"`
	Author         string         `json:"author"`
	Category       sql.NullString `json:"category"`
	Content        []byte         `json:"content"`
	CoverImage     []byte         `json:"cover_image"`
	CoverImageType sql.NullString `json:"cover_image_type"`
	ID             int32          `json:"id"`
}

func (q *Queries) UpdateBook(ctx context.Context, arg UpdateBookParams) error {
	_, err := q.exec(ctx, q.updateBookStmt, updateBook,
		arg.Title,
		arg.Author,
		arg.Category,
		arg.Content,
		arg.CoverImage,
		arg.CoverImageType,
		arg.ID,
	)
	return err
}
