-- name: GetBook :one
SELECT
    id,
    title,
    author,
    category,
    content,
    cover_image,
    cover_image_type
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT
    id,
    title,
    author,
    category,
    cover_image_type
FROM books
ORDER BY id;

-- name: CreateBook :exec
INSERT INTO books (
    title,
    author,
    category,
    content,
    cover_image,
    cover_image_type
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateBook :exec
UPDATE books
SET
    title = $1,
    author = $2,
    category = $3,
    content = $4,
    cover_image = $5,
    cover_image_type = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $7;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = $1;

-- name: SearchBooksByTitle :many
SELECT
    id,
    title,
    author,
    category,
    cover_image_type
FROM books
WHERE title LIKE '%' || $1 || '%'
ORDER BY title;

-- name: SearchBooksByAuthor :many
SELECT
    id,
    title,
    author,
    category,
    cover_image_type
FROM books
WHERE author LIKE '%' || $1 || '%'
ORDER BY author;

-- name: CountBooks :one
SELECT COUNT(*) AS count FROM books;
