CREATE TABLE books (
    id SERIAL PRIMARY KEY,                        -- Auto-incrementing integer ID
    title TEXT NOT NULL,                          -- Title of the book (required)
    author TEXT NOT NULL,                         -- Author of the book (required)
    category TEXT,                                -- Optional category of the book
    content BYTEA,                                -- Binary data for book content (PDF)
    cover_image BYTEA,                            -- Binary data for the cover image
    cover_image_type TEXT,                        -- MIME type of the cover image (e.g., "image/png")
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp of last update
);

CREATE UNIQUE INDEX idx_books_title_author ON books (title, author);
